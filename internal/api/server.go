package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/roky61980-create/port-monitor/internal/alerts"
	"github.com/roky61980-create/port-monitor/internal/events"
	"github.com/roky61980-create/port-monitor/internal/model"
	"github.com/roky61980-create/port-monitor/internal/store"
	"github.com/roky61980-create/port-monitor/internal/telemetry"
)

type Server struct { db *store.Store; hub *events.Hub; alerts *alerts.Engine; otel *telemetry.Telemetry }
func New(db *store.Store,h *events.Hub,a *alerts.Engine,t *telemetry.Telemetry)*Server{return &Server{db:db,hub:h,alerts:a,otel:t}}
func (s *Server) Handler()http.Handler{m:=http.NewServeMux();m.HandleFunc("/api/health",func(w http.ResponseWriter,r *http.Request){jsonOut(w,200,map[string]string{"status":"ok"})});m.HandleFunc("/api/metrics",func(w http.ResponseWriter,r *http.Request){jsonOut(w,200,map[string]string{"service":"port-monitor","telemetry":"otlp-traces"})});m.HandleFunc("/api/devices",s.devices);m.HandleFunc("/api/devices/",s.device);m.HandleFunc("/api/ports",s.ports);m.HandleFunc("/api/optical",s.optical);m.HandleFunc("/api/events",s.eventList);m.HandleFunc("/api/alerts",s.alertList);m.HandleFunc("/api/alerts/",s.alertAction);m.HandleFunc("/api/topology",s.topology);m.HandleFunc("/api/events/stream",s.stream);m.Handle("/",http.FileServer(http.Dir("web")));return m}
func (s *Server) devices(w http.ResponseWriter,r *http.Request){switch r.Method{case http.MethodGet:x,e:=s.db.Devices();if e!=nil{fail(w,e);return};jsonOut(w,200,x);case http.MethodPost:var d model.Device;if !decode(w,r,&d){return};if d.ID==""||d.IP==""{jsonOut(w,400,map[string]string{"error":"id and ip are required"});return};if d.Status==""{d.Status="UNKNOWN"};d.Enabled=true;if e:=s.db.UpsertDevice(d);e!=nil{fail(w,e);return};jsonOut(w,201,d);default:w.WriteHeader(405)}}
func (s *Server) device(w http.ResponseWriter,r *http.Request){p:=strings.Trim(strings.TrimPrefix(r.URL.Path,"/api/devices/"),"/");parts:=strings.Split(p,"/");if len(parts)>1&&parts[1]=="ports"{x,e:=s.db.Ports();if e!=nil{fail(w,e);return};out:=[]model.Port{};for _,v:=range x{if v.DeviceID==parts[0]{out=append(out,v)}};jsonOut(w,200,out);return};if len(parts)>1&&parts[1]=="optical"{x,e:=s.db.Optical();if e!=nil{fail(w,e);return};out:=[]model.OpticalReading{};for _,v:=range x{if v.DeviceID==parts[0]{out=append(out,v)}};jsonOut(w,200,out);return};d,e:=s.db.Device(parts[0]);if e!=nil{jsonOut(w,404,map[string]string{"error":"device not found"});return};switch r.Method{case http.MethodGet:jsonOut(w,200,d);case http.MethodPatch:var patch model.Device;if !decode(w,r,&patch){return};if patch.Name!=""{d.Name=patch.Name};if patch.Model!=""{d.Model=patch.Model};if patch.Vendor!=""{d.Vendor=patch.Vendor};if patch.Location!=""{d.Location=patch.Location};if patch.Description!=""{d.Description=patch.Description};d.Enabled=patch.Enabled;if e=s.db.UpsertDevice(d);e!=nil{fail(w,e);return};jsonOut(w,200,d);case http.MethodDelete:if e=s.db.DeleteDevice(d.ID);e!=nil{fail(w,e);return};w.WriteHeader(204);default:w.WriteHeader(405)}}
func (s *Server) ports(w http.ResponseWriter,r *http.Request){x,e:=s.db.Ports();if e!=nil{fail(w,e);return};jsonOut(w,200,x)}
func (s *Server) optical(w http.ResponseWriter,r *http.Request){x,e:=s.db.Optical();if e!=nil{fail(w,e);return};jsonOut(w,200,x)}
func (s *Server) eventList(w http.ResponseWriter,r *http.Request){limit,_:=strconv.Atoi(r.URL.Query().Get("limit"));offset,_:=strconv.Atoi(r.URL.Query().Get("offset"));if limit<1||limit>500{limit=100};if offset<0{offset=0};x,e:=s.db.Events(limit,offset);if e!=nil{fail(w,e);return};jsonOut(w,200,x)}
func (s *Server) alertList(w http.ResponseWriter,r *http.Request){x,e:=s.db.Alerts();if e!=nil{fail(w,e);return};jsonOut(w,200,x)}
func (s *Server) alertAction(w http.ResponseWriter,r *http.Request){p:=strings.Trim(strings.TrimPrefix(r.URL.Path,"/api/alerts/"),"/");parts:=strings.Split(p,"/");a,e:=s.db.Alert(parts[0]);if e!=nil{jsonOut(w,404,map[string]string{"error":"alert not found"});return};if r.Method!=http.MethodPost||len(parts)!=2{w.WriteHeader(405);return};now:=time.Now().UTC();switch parts[1]{case "ack":a.Status="ACKNOWLEDGED";a.AcknowledgedBy=r.Header.Get("X-Actor");case "resolve":a.Status="RESOLVED";a.ResolvedAt=&now;default:http.NotFound(w,r);return};a.UpdatedAt=now;if e=s.db.SaveAlert(a);e!=nil{fail(w,e);return};s.hub.Publish(map[string]any{"kind":"alert","event_type":"ALERT_"+strings.ToUpper(parts[1])+"ED","alert":a});if s.otel!=nil{s.otel.Event(r.Context(),"ALERT_"+strings.ToUpper(parts[1]),a.Device,"","",a.Interface,a.Severity,a.Status)};jsonOut(w,200,a)}
func (s *Server) topology(w http.ResponseWriter,r *http.Request){d,e:=s.db.Devices();if e!=nil{fail(w,e);return};p,e:=s.db.Ports();if e!=nil{fail(w,e);return};jsonOut(w,200,map[string]any{"devices":d,"ports":p})}
func (s *Server) stream(w http.ResponseWriter,r *http.Request){f,ok:=w.(http.Flusher);if !ok{http.Error(w,"streaming unsupported",500);return};w.Header().Set("Content-Type","text/event-stream");w.Header().Set("Cache-Control","no-cache");w.Header().Set("Connection","keep-alive");ch,cancel:=s.hub.Subscribe();defer cancel();keep:=time.NewTicker(20*time.Second);defer keep.Stop();for{select{case <-r.Context().Done():return;case b,ok:=<-ch:if !ok{return};_,_=w.Write(append(append([]byte("data: "),b...),[]byte("\n\n")...));f.Flush();case <-keep.C:_,_=w.Write([]byte(": keepalive\n\n"));f.Flush()}}}
func decode(w http.ResponseWriter,r *http.Request,v any)bool{if e:=json.NewDecoder(http.MaxBytesReader(w,r.Body,1<<20)).Decode(v);e!=nil{jsonOut(w,400,map[string]string{"error":"invalid JSON"});return false};return true}
func jsonOut(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func fail(w http.ResponseWriter,e error){if errors.Is(e,context.Canceled){jsonOut(w,503,map[string]string{"error":"request canceled"});return};jsonOut(w,500,map[string]string{"error":"internal server error"})}

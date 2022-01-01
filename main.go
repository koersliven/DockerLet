package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"net"
	"os"
	"flag"
	"github.com/mindprince/gonvml"
	"log"
	"bytes"
	"net/http"
	"os/exec"
	"math"
	"syscall"
	"io/ioutil"
	"strings"
	"path/filepath"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
)

type IdentifyInfo struct{
	RackID	 		uint 			`json:"rackID"`	
	HostName    	string 			`json:"hostName"`
	HostIP    		string 			`json:"hostIP"`
	DataPathIP    	string 			`json:"dataPathIP"`
	State	 		string 			`json:"state"`
	XpuTypeNum	 	int 			`json:"xpuTypeNum"`
	TotalVxpuNum	 	int 			`json:"totalVxpuNum"`
}

type AttitudeInfo struct{
	XpuName    		string 			`json:"xpuName"`
	Version    		string 			`json:"version"`
	IfType    		string 			`json:"ifType"`
	TotalMemory         	string 		    	`json:"totalMemory"`
	TotalInt4	 	string 		    	`json:"totalInt4"`
	MaxInt4	 	        string 		    	`json:"maxInt4"`
	MinInt4	 	        string 		    	`json:"minInt4"`
	TotalInt8 	        string 		    	`json:"totalInt8"`
	MaxInt8 	        string 		    	`json:"maxInt8"`
	MinInt8	 	        string 		    	`json:"minInt8"`
	TotalFP16 	        string 		    	`json:"totalFP16"`
	MaxFP16 	        string 		    	`json:"maxFP16"`
	MinFP16	 	        string 		    	`json:"minFP16"`
	TotalFP32 	        string 		    	`json:"totalFP32"`
	MaxFP32 	        string 		    	`json:"maxFP32"`
	MinFP32	 	        string 		    	`json:"minFP32"`
	Fp32ConvInt4 	    	string 		    `json:"fp32ConvInt4"`
	Fp32ConvInt8	 	string 		    	`json:"fp32ConvInt8"`
	Fp32ConvFp16	 	string 		    	`json:"fp32ConvFp16"`
	Fp16ConvInt4 	    	string 		    `json:"fp16ConvInt4"`
	Fp16ConvInt8	 	string 		    	`json:"fp16ConvInt8"`
	Fp16ConvFp32	 	string 		    	`json:"fp16ConvFp32"`
	Int8ConvInt4	 	string 		    	`json:"int8ConvInt4"`
	Int8ConvFp16	 	string 		    	`json:"int8ConvFp16"`
	Int8ConvFp32	 	string 		    	`json:"int8ConvFp32"`
	Int4ConvInt8	 	string 		    	`json:"int4ConvInt8"`
	Int4ConvFp16	 	string 		    	`json:"int4ConvFp16"`
	Int4ConvFp32	 	string 		    	`json:"int4ConvFp32"`
	XpuNum	 		uint 		    		`json:"xpuNum"`
}

type Xpu struct {
    XpuName    		string 			`json:"xpuName"`
   	XpuUUID    		string 			`json:"xpuUUID"`
	XpuDeviceID    		int 			`json:"xpuDeviceID"`
    State		       	string    		`json:"state"`
    LeftMemory  	    	string 			`json:"leftMemory"`
	LeftInt4         	string 		    	`json:"leftInt4"`
	LeftInt8		string 		    	`json:"leftInt8"`
	LeftFp16		string 		    	`json:"leftFp16"`
	LeftFp32		string 		    	`json:"leftFp32"`

	VxpuNum	 		uint 		    	`json:"vxpuNum"`
    Vxpus	     		[]Vxpu   		`json:"vxpus"`
}	

type Vxpu struct {
    XpuUUID    		string 			`json:"xpuUUID"`
    State		       	string    		`json:"state"`
    AllocatedMemory  	string 			`json:"allocatedMemory"`
	AllocateInt4        	string 		    	`json:"allocateInt4"`
	AllocateInt8		string 		    	`json:"allocateInt8"`
	AllocateFp16		string 		    	`json:"allocateFp16"`
	AllocateFp32		string 		    	`json:"allocateFp32"`
    ContainerName           string    	        `json:"containerName"`
    ContainerPort   	string 		        `json:"containerPort"`
	ContainerIP			string				`json:"containerIP"`
}	

type Metrics struct{
	HnIdentifyInfo		IdentifyInfo   		`json:"hnIdentifyInfo"`
	XpuAttitude		[]AttitudeInfo   	`json:"xpuAttitude"`
	Xpus	     		[]Xpu   		`json:"xpus"`
}

type ExtendInfo struct {
	XpuName    		string 			`json:"xpuName"`
	XpuUUID    		string 			`json:"xpuUUID"`
}

type AllocReq2 struct {
	//HN info for double check
	
	ApplyType 		 string     `json:"applyType"`
	Extendinfo	    map[string]string   	`json:"extendInfo"`
	Flops			 string 		`json:"flops"`
	Size			string		`json:"size"`
	Memory         	string   	`json:"memory"`
	Precision    	string 		`json:"precision"`
	TransId			string		`json:"transId"`
}

type AllocReq3 struct {
	ApplyList	    []AllocReq2   	`json:"applyList"`
}

type AllocResp struct {
	State					string				   `json:"success"` 
    Result                  string                  `json:"errorMsg"`
	VxpuInfo	       	[]Allocresplist      	`json:"applyResponsesList"`
	AfterAssign			[]AfterAssignXpu 		`json:"afterAssign"`
}

type Allocresplist struct {
	TransId			string		`json:"transId"`
	VxpuUUID    		string 			`json:"VxpuUUID"`
	ExtendInfo	    map[string]string   	`json:"extendInfo"`
}

type AllocrespExtendInfo struct {
	ContainerName			string		`json:"containerName"`
	ContainerPort    		string 			`json:"containerPort"`
	ContainerIP	    string   	`json:"containerIP"`
}

type AfterAssignXpu struct {
    XpuName    		string 			`json:"xpuName"`
   	XpuUUID    		string 			`json:"xpuUUID"`
	XpuDeviceID    		int 			`json:"xpuDeviceID"`
    State		       	string    		`json:"state"`
    LeftMemory  	    	string 			`json:"leftMemory"`
	LeftInt4         	string 		    	`json:"leftInt4"`
	LeftInt8		string 		    	`json:"leftInt8"`
	LeftFp16		string 		    	`json:"leftFp16"`
	LeftFp32		string 		    	`json:"leftFp32"`
	VxpuNum	 		uint 		    	`json:"vxpuNum"`
}

type ReleaseReq struct {
	//Vxpu info
	ReleaseList	       	[]Releaselist   	`json:"releaseList"` 
}

type Releaselist struct {
	//Vxpu info
	VxpuUUID	       	string   	`json:"vxpuUUID"`
	ExtendInfo		map[string]string	`json:"extendInfo"`
}

type Releaseresplist struct {
	//Vxpu info
	VxpuUUID	    string   	`json:"vxpuUUID"`
	State		string	 `json:"success"`
	Extendinfo	    map[string]string   	`json:"extendInfo"`
}

type ReleaseResp struct {
	State					string			 `json:"success"`
    Result                  string           `json:"errorMsg"` 
	VxpuInfo	       	[]Releaseresplist   	`json:"releaseList"`
	AfterAssign			[]AfterAssignXpu 		`json:"afterAssign"`
}

type AllocBackup struct {
    Index                   int
    LeftMemory  	    	string
	LeftInt4         string
	LeftInt8         string
	LeftFp16         string
	LeftFp32         string
	Percentage		 string
    VxpuInfo                Vxpu
	Memory 					string
	Flops 					string
	Precision 				string
}

type Testmap struct {
	//Vxpu info
	Mymap		map[string]string	 `json:"extendInfo"`
}

const(
	version="0.1.0"
	flavor = "runtime"
	base_image = "cuda10.2-cudnn8.1.1.33-trt8.0.3-ubuntu18.04"
	name_space = "reg.docker.alibaba-inc.com/vodla_ais"
	g_data_net_name = "ens5f1"
	g_host_net_name = "bond0"
	g_image = name_space + "/vxpu" + ":" + version + "-" + flavor + "-" + base_image
	g_container_port_start = 30100
	g_port = "9405"
)
// container port start number
var port_start int64

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func init() {
    // 以JSON格式为输出，代替默认的ASCII格式
    //logrus.SetFormatter(&logrus.JSONFormatter{})
    // 以Stdout为输出，代替默认的stderr
    //logrus.SetOutput(os.Stdout)
    // 设置日志等级
    //logrus.SetLevel(logrus.WarnLevel)
	logrus.SetLevel(logrus.InfoLevel)
	// 设置在输出日志中添加文件名和方法信息
	logrus.SetReportCaller(true)
}

//----------------start-----------------//

func main(){
	// === 命令行参数解析 ===
	// vxpu container image
	var image string
	// the path of tmp file
	var tmp_file_path string
	// the data net
	var data_net_name string
	// the host net
	var host_net_name string
	// server port number
	var port string
	// container name
	var c_name string
	// xpu_info
	var xpu_info_file string

	home, err :=homeUnix()

	flag.StringVar(&image, "i", g_image, "container image")
	flag.StringVar(&tmp_file_path, "d", home + "/.vxpulet", "tmp file path")
	flag.StringVar(&xpu_info_file, "f", home + "/.vxpulet/xpuinfo.json", "xpu info file")
	flag.StringVar(&port, "p", g_port, "RESTful port")
	flag.StringVar(&c_name, "u", "vxpu_", "container name")
	flag.StringVar(&data_net_name, "c", g_data_net_name, "container net name")
	flag.StringVar(&host_net_name, "h", g_host_net_name, "host net name")
	flag.Int64Var(&port_start, "ps", g_container_port_start, "container start port")

	// enable 上面的命令行
	flag.Parse()

	// 设置日志重定向
	writer1 := &bytes.Buffer{}
	//writer2 := os.Stdout
	Shellout("mkdir -p "+ tmp_file_path)
	writer3, err := os.OpenFile(tmp_file_path+"/log.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
	  log.Fatalf("create file log.txt failed: %v", err)
	}
  
	//logrus.SetOutput(io.MultiWriter(writer1, writer2, writer3))
	logrus.SetOutput(io.MultiWriter(writer1, writer3))

	// 打印参数信息
	logrus.Info(">>>>>> start args >>>>>>>")
	logrus.Info("image: ", image)
	logrus.Info("tmp path: ", tmp_file_path)
	logrus.Info("xpu info: ", xpu_info_file)
	logrus.Info("port: ", port)
	logrus.Info("name: ", c_name)
	logrus.Info("data net: ", data_net_name)
	logrus.Info("host net: ", host_net_name)
	logrus.Info("start port: ", port_start)
	logrus.Info("<<<<<<< start args <<<<<<<")

	// 资源上报参数
	var metrics Metrics
	
	var ips IpInfo
	var hostIP string
	var dataPathIP IpInfo 
	localAddresses(&ips)
	for i := 0; i < 2; i++ {
		logrus.Info("name: ", ips.name[i], " ,addr: ", ips.addr[i])
		if ips.name[i] == g_host_net_name {
			hostIP = ips.addr[i]
		} else if ips.name[i] == g_data_net_name {
			dataPathIP.name[0] = ips.name[i]
			dataPathIP.addr[0] = ips.addr[i]
		}
	}

	// Check if xpuinfo.json file exist
	// 如果xpuinfo.json不存在，则生成新的设备信息文件
    b, err := PathExists(xpu_info_file)
    if err != nil {
        logrus.Error("PathExists(%s),err(%v)\n", xpu_info_file, err)
    }
    if b {
        filePtr, err := os.Open(xpu_info_file)
        if err != nil {
            logrus.Error("文件打开失败 [Err:%s]", err.Error())
            return
        }
		logrus.Trace("readed "+xpu_info_file)
        defer filePtr.Close()
        // 创建json解码器
        decoder := json.NewDecoder(filePtr)
        err = decoder.Decode(&metrics)
        if err != nil {
            logrus.Error("metrics 解码失败", err.Error())
        } else {
            logrus.Trace("metrics 解码成功")
            logrus.Trace(metrics)
        }
        filePtr.Close()
    } else {
        logrus.Trace("xpu_info_file %s not exist\n", xpu_info_file)
		var xpu_info_dir string
		xpu_info_dir = getParentDirectory(xpu_info_file);
		Shellout("mkdir -p "+ xpu_info_dir)

        //If not exist: create and add HN/XPU info, vXpu num is all zero
        if err := gonvml.Initialize(); err != nil {
            logrus.Fatal("gonvml.Initialize() failed.")
        }
        defer gonvml.Shutdown()
        
        version, err := gonvml.SystemDriverVersion()
        if err != nil {
            logrus.Fatal("gonvml.SystemDriverVersion() failed.")
        }
        
        numDevices, err := gonvml.DeviceCount()
        if err != nil {
            logrus.Fatal("gonvml.DeviceCount() failed.")
        }
        
        hostname, err := os.Hostname()
        if err != nil {
        	logrus.Fatal("os.Hostname() failed.")
        }
          
		// metrics资源查询
        xpuTypeNum := 1

        identify := IdentifyInfo{
        	RackID: 0,
        	HostName: hostname,
			HostIP: hostIP,
			DataPathIP: dataPathIP.addr[0],
        	State: "Avaliable",
        	XpuTypeNum: xpuTypeNum,
        }
    
        var attitude []AttitudeInfo
        for typenum := 0; typenum < xpuTypeNum; typenum++ {
            attitude = append(attitude,
				// 算力的种类和数据要可以自动获取，不同机器不同，这里要大改
           	    AttitudeInfo{
           		Version: version,
           		IfType: "PCIe",
           		TotalInt4: "260TOPS",
           		MaxInt4: "260TOPS",
           		MinInt4: "1TOPS",
           		TotalInt8: "130TOPS",
           		MaxInt8: "130TOPS",
           		MinInt8: "1TOPS",
           		TotalFP16: "65TF",
           		MaxFP16: "65TF",
           		MinFP16: "1TF",
           		TotalFP32: "8.1TF",
           		MaxFP32: "8.1TF",
           		MinFP32: "0.1TF",
           		Fp32ConvInt4: "8",
           		Fp32ConvInt8: "4",
           		Fp32ConvFp16: "2",
           		Fp16ConvInt4: "4",
           		Fp16ConvInt8: "2",
           		Fp16ConvFp32: "2",
           		Int8ConvInt4: "2",
           		Int8ConvFp16: "0.5",
           		Int8ConvFp32: "0.0625",
           		Int4ConvInt8: "0.4",
           		Int4ConvFp16: "0.25",
           		Int4ConvFp32: "0.3125",
           		XpuNum: numDevices,
            })
        }
        
        metrics.HnIdentifyInfo = identify
        metrics.XpuAttitude = attitude
    
        for index := 0; index < int(numDevices); index++ {
            device, err := gonvml.DeviceHandleByIndex(uint(index))
            if err != nil {
                logrus.Fatal("gonvml.DeviceHandleByIndex failed.")
            }
    
            gpuname, err := device.Name()
            if err != nil {
            	logrus.Fatal("device.Name() failed.")
            }
            metrics.XpuAttitude[0].XpuName = gpuname
            
            memoryTotal, _, err := device.MemoryInfo()
            if err != nil {
            	logrus.Fatal("device.MemoryInfo() failed.")
            }
            metrics.XpuAttitude[0].TotalMemory = strconv.Itoa(int(memoryTotal)/1024/1024)+"MB"
            
            gpuuuid, err := device.UUID()
            if err != nil {
                logrus.Fatal("device.UUID() failed.")
            }
            
            var xpu Xpu
            xpu.XpuName = gpuname
            xpu.XpuUUID =               gpuuuid
			xpu.XpuDeviceID = 			index
            xpu.State =                 "Avaliable"
            xpu.LeftMemory =            metrics.XpuAttitude[0].TotalMemory
            xpu.LeftInt4 =              metrics.XpuAttitude[0].TotalInt4
            xpu.LeftInt8 =              metrics.XpuAttitude[0].TotalInt8
            xpu.LeftFp16 =              metrics.XpuAttitude[0].TotalFP16
            xpu.LeftFp32 =              metrics.XpuAttitude[0].TotalFP32
            xpu.VxpuNum =               0
            metrics.Xpus = append(metrics.Xpus, xpu)  //增添Xpu0分配的vXpus
        }
	
		// 创建文件
		filePtr, err := os.Create(xpu_info_file) 
		if err != nil {
			logrus.Error("文件创建失败", err.Error())
			return
		}
		logrus.Trace("created "+xpu_info_file)
		defer filePtr.Close()

		// // 创建Json编码器
		// encoder := json.NewEncoder(filePtr)
		// err = encoder.Encode(metrics)

		encoder, err := json.Marshal(metrics)
		err = ioutil.WriteFile(xpu_info_file, encoder, 0777)

		if err != nil {
			logrus.Error("编码错误", err.Error())
		} else {
			logrus.Trace("编码成功")
		}
		filePtr.Close()
    }
	// 以上，xpuinfo.json的基本读写完成
	
	// metrics-server
	// 用来读取 xpuinfo.json 的信息，并返回
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// 根据w创建Encoder，然后调用Encode()方法将对象编码成JSON
		// Encoder 主要负责将结构对象编码成 JSON 数据
		// 调用json.NewEncoder(io.Writer)方法获得一个 Encoder 实例
		logrus.Info("====== metrics =======")
		var metrics_load Metrics
        fp_metrics, err := os.Open(xpu_info_file)
        if err != nil {
            logrus.Error("文件打开失败 [Err:%s]", err.Error())
            return
        }
        // defer fp_metrics.Close()

		// 阻塞模式下加排他锁
		if err := syscall.Flock(int(fp_metrics.Fd()), syscall.LOCK_EX); err != nil {
			logrus.Error("metrics上锁失败", err)
		}else{
			logrus.Trace("metrics上锁成功")
		}

		logrus.Trace("-----------------开始读取当前设备信息---------------------")
        // 创建json解码器
        decoder := json.NewDecoder(fp_metrics)
        err = decoder.Decode(&metrics_load)

        if err != nil {
            logrus.Error("metrics解码失败", err.Error())
        } else {
            logrus.Trace("metrics解码成功")
        }

		// metrics查询信息 
		json.NewEncoder(w).Encode(metrics_load)

		// 解锁
		if err := syscall.Flock(int(fp_metrics.Fd()), syscall.LOCK_UN); err != nil {
			logrus.Error("metrics解锁失败", err)
		}else{
			logrus.Trace("metrics解锁成功")
		}

        fp_metrics.Close()
    })

	// 分配资源
	http.HandleFunc("/allocate", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("======== allocate =========")
		// 从文件获取当前设备信息并且给文件上锁
		var metrics_allocate Metrics
		fp_allocate, err := os.OpenFile(xpu_info_file, os.O_RDWR, os.ModePerm)

		if err != nil {
			logrus.Error("文件打开失败 [Err:%s]", err.Error())
			return
		}

		// 阻塞模式下加排他锁
		if err := syscall.Flock(int(fp_allocate.Fd()), syscall.LOCK_EX); err != nil {
			logrus.Error("allocate上锁失败", err)
		}else{
			logrus.Trace("allocate上锁成功")
		}

		// defer fp_allocate.Close()
		// 创建json解码器
		decoder := json.NewDecoder(fp_allocate)
		err = decoder.Decode(&metrics_allocate)

		logrus.Trace("-----------------开始分配进程---------------------")

		// var allocreq2 AllocReq2  // input
		// 读取接口的申请 request
		var allocreq3 AllocReq3  // input request
		var allocresp AllocResp // output response

		// 解码器
		json.NewDecoder(r.Body).Decode(&allocreq3)
	
		// 记录可以正确分配的容器个数
		successcnt := 0 // 从设备信息检查成功分配的个数
		finalcnt := 0 // 检查docker生成的成功个数
		// 记录请求里面要生成几个容器
		paralen := len(allocreq3.ApplyList)
		// 上面的三个参数相等时，才会开始创造容器
		
		// 记录当前请求下有哪些xpu改变，注意是xpu不是vxpu
		var changexpu []string
	
		// 记录当前要生成哪些容器的编号
		var create_cnt_num []int
		var oldcnt_max int
	
		// metrics更新前备用
		var backuparray []AllocBackup
	
		// 创建临时变量，为了检查有无报错，但是不会动到xpuinfo.json
		var tmp_metrics Metrics
		for i := 0; i < len(metrics_allocate.Xpus); i++ {
			tmp_metrics.Xpus = append(tmp_metrics.Xpus, metrics_allocate.Xpus[i])
		}
	
		allocresp.State = "false"  // "success": "false",
	
		// 检查入参
		if paralen == 0 {
			allocresp.Result = "applyList in empty"  // "errorMsg": "applyList in empty",
			goto CONSTRUCT_RESP  // 结束，将metrics_allocate返回至xpuinfo.json	
		}
		
		// 检查每个请求的入参
		// 针对 byFlops。
		// 如果是bySize，这里要改
		for nodenum := 0; nodenum < paralen; nodenum++ {
			if (allocreq3.ApplyList[nodenum].Precision == "") || 
				(allocreq3.ApplyList[nodenum].Memory == "") ||
				(allocreq3.ApplyList[nodenum].Flops == "") ||
				(allocreq3.ApplyList[nodenum].Extendinfo["xpuName"] == "") ||
				(allocreq3.ApplyList[nodenum].Extendinfo["xpuUUID"] == ""){
					allocresp.Result = fmt.Sprintf("Parameter uncorrect in applyList[%d]", nodenum)
					goto CONSTRUCT_RESP
			}				
		}
	
		// 开始分配内存和算力
		for nodenum := 0; nodenum < paralen; nodenum++ {
			// 未扩展 bySize 功能
			if allocreq3.ApplyList[nodenum].ApplyType == "bySize" {
				allocresp.Result = fmt.Sprintf("ApplyType unsupport in applyList[%d]", nodenum)
				goto CONSTRUCT_RESP
			}else {

				// 去掉内存memory的单位
				allocreq3_mem := allocreq3.ApplyList[nodenum].Memory[:len(allocreq3.ApplyList[nodenum].Memory)-2]

				// 去掉算力flops的单位
				var allocreq3_Flops string
				if allocreq3.ApplyList[nodenum].Precision == "Int4" ||
					allocreq3.ApplyList[nodenum].Precision == "Int8" {
						// allocreq3.ApplyList[nodenum].Flops = allocreq3.ApplyList[nodenum].Flops[:len(allocreq3.ApplyList[nodenum].Flops)-4]
						allocreq3_Flops = allocreq3.ApplyList[nodenum].Flops[:len(allocreq3.ApplyList[nodenum].Flops)-4]
						} else {
						// allocreq3.ApplyList[nodenum].Flops = allocreq3.ApplyList[nodenum].Flops[:len(allocreq3.ApplyList[nodenum].Flops)-2]
						allocreq3_Flops = allocreq3.ApplyList[nodenum].Flops[:len(allocreq3.ApplyList[nodenum].Flops)-2]
				}

				for typenum := 0; typenum < metrics_allocate.HnIdentifyInfo.XpuTypeNum; typenum++ {
					// 获取总算力
					totalint4cpt := metrics_allocate.XpuAttitude[typenum].TotalInt4 // 当前xpu的Int4的总算力
					totalint4cptdata,_ := strconv.ParseFloat(totalint4cpt[:len(totalint4cpt)-4], 32)
					totalint8cpt := metrics_allocate.XpuAttitude[typenum].TotalInt8 // 
					totalint8cptdata,_ := strconv.ParseFloat(totalint8cpt[:len(totalint8cpt)-4], 32)
					totalfp16cpt := metrics_allocate.XpuAttitude[typenum].TotalFP16 // 
					totalfp16cptdata,_ := strconv.ParseFloat(totalfp16cpt[:len(totalfp16cpt)-2], 32)
					totalfp32cpt := metrics_allocate.XpuAttitude[typenum].TotalFP32 // 
					totalfp32cptdata,_ := strconv.ParseFloat(totalfp32cpt[:len(totalfp32cpt)-2], 32)
					totalfp32cptdata = remainfp32(totalfp32cptdata)
	
					if allocreq3.ApplyList[nodenum].Extendinfo["xpuName"] == metrics_allocate.XpuAttitude[typenum].XpuName {
						index := 0 
	
						// 针对不同xpu group的最小切割算力
						// var flop_mini float64 = 5
						var flop_mini float64
						flop_mini = cut_minisize(metrics_allocate.XpuAttitude[typenum].XpuName)
						if flop_mini == 0 {
							allocresp.Result = fmt.Sprintf("No satisfied flops_cut_minisize in applyList[%d]", nodenum)
							goto CONSTRUCT_RESP
						}

						// 循环同一类型的GPU
						for ; index < int(metrics_allocate.XpuAttitude[typenum].XpuNum); index++ {
							// 对上UUID，我要的在HN有没有
							if allocreq3.ApplyList[nodenum].Extendinfo["xpuUUID"] == metrics_allocate.Xpus[index].XpuUUID {
								var backup AllocBackup
								var percentcpt float64
								backup.Percentage = "0.0%"

								//check if can allocate
								// 检查剩的还能不能分配
								leftmem := tmp_metrics.Xpus[index].LeftMemory // "LeftMemory":"15109MB"
								memdata,_ := strconv.Atoi( leftmem[:len(leftmem)-2] ) // 有单位，剪掉 
								reqmem,_ := strconv.ParseFloat(allocreq3_mem, 32)
								// 需求超过了现存，分不了
								if memdata < (int)(reqmem) {
									allocresp.Result = fmt.Sprintf("Memory over allocate in applyList[%d]", nodenum)
									goto CONSTRUCT_RESP
								} else {
									backup.VxpuInfo.AllocatedMemory = allocreq3.ApplyList[nodenum].Memory
									// backup.Memory = allocreq3.ApplyList[nodenum].Memory
									// step0. 预处理，求得目前四种算力还剩的分配空间（同步分配所以四个都要算）
									leftcptInt4 := tmp_metrics.Xpus[index].LeftInt4 // 
									cptdataInt4,_ := strconv.ParseFloat(leftcptInt4[:len(leftcptInt4)-4], 32) // Int4
									leftcptInt8 := tmp_metrics.Xpus[index].LeftInt8
									cptdataInt8,_ := strconv.ParseFloat(leftcptInt8[:len(leftcptInt8)-4], 32) // Int8
									leftcptFp16 := tmp_metrics.Xpus[index].LeftFp16
									cptdataFp16,_ := strconv.ParseFloat(leftcptFp16[:len(leftcptFp16)-2], 32) // Fp16
									leftcptFp32 := tmp_metrics.Xpus[index].LeftFp32
									cptdataFp32,_ := strconv.ParseFloat(leftcptFp32[:len(leftcptFp32)-2], 32) // Fp32
									cptdataFp32 = remainfp32(cptdataFp32)
	
									// step1. 不管啥精度，先得到请求的算力分配，但不是最终确定分配的算力资源
									reqcpt,_ := strconv.ParseFloat(allocreq3_Flops, 32) 
									log.Println("第",nodenum,"个vxpu，此时请求的算力是",allocreq3.ApplyList[nodenum].Precision,"的", reqcpt)
	
									// step2. 四个elseif只求算力分配百分比，t4最小粒度是5%
									if allocreq3.ApplyList[nodenum].Precision == "Int4" {
										percentcpt = cpt_percent_t4(flop_mini, reqcpt, totalint4cptdata)
										// log.Println("第",nodenum,"个vxpu，此时应该被分配的算力百分比是", percentcpt)	
									}else if allocreq3.ApplyList[nodenum].Precision == "Int8" {
										percentcpt = cpt_percent_t4(flop_mini, reqcpt, totalint8cptdata)
										// log.Println("第",nodenum,"个vxpu，此时应该被分配的算力百分比是", percentcpt)	
									} else if allocreq3.ApplyList[nodenum].Precision == "Fp16" {
										percentcpt = cpt_percent_t4(flop_mini, reqcpt, totalfp16cptdata)
										// log.Println("第",nodenum,"个vxpu，此时应该被分配的算力百分比是", percentcpt)	
									} else if allocreq3.ApplyList[nodenum].Precision == "Fp32" {
										percentcpt = cpt_percent_t4(flop_mini, reqcpt, totalfp32cptdata)
										// log.Println("第",nodenum,"个vxpu，此时应该被分配的算力百分比是", percentcpt)	
									} else {
										// 算力不支持，没有百分比了
										// allocresp.Result = "Compute unsupport"
										allocresp.Result = fmt.Sprintf("ComputeType unsupport in applyList[%d]", nodenum)
										goto CONSTRUCT_RESP
									}
									
									backup.Percentage = strconv.FormatFloat(percentcpt, 'f', -1, 32)+"%"
	
									// step3. 根据百分比计算实际要分配出去的算力资源（同步分配所以四个都要算）
									realcutcptInt4 := totalint4cptdata * percentcpt / 100.0 //TOPS
									realcutcptInt8 := totalint8cptdata * percentcpt / 100.0 //TOPS
									realcutcptFp16 := totalfp16cptdata * percentcpt / 100.0 //TF
									realcutcptFp32 := remainfp32(totalfp32cptdata * percentcpt / 100.0) //TF
									log.Println("第",nodenum,"个vxpu，真实分配的算力都是", realcutcptInt4, realcutcptInt8, realcutcptFp16, realcutcptFp32)
	
									// step4. 再计算要分配的算力有无超出现存剩余
									if (cptdataInt4 < realcutcptInt4) ||
										(cptdataInt8 < realcutcptInt8) ||
										(cptdataFp16 < realcutcptFp16) ||
										(cptdataFp32 < realcutcptFp32) {
										allocresp.Result = fmt.Sprintf("Compute over allocate in applyList[%d]", nodenum)
										goto CONSTRUCT_RESP
									}
									
									// step5. 更新临时变量backup，为了后面的metrics更新
									backup.Index = index
									backup.LeftMemory = strconv.Itoa( int( memdata - (int)(reqmem) ) )+"MB"
									backup.LeftInt4 = strconv.FormatFloat(cptdataInt4 - realcutcptInt4, 'f', -1, 32)+"TOPS"
									backup.LeftInt8 = strconv.FormatFloat(cptdataInt8 - realcutcptInt8, 'f', -1, 32)+"TOPS"
									backup.LeftFp16 = strconv.FormatFloat(cptdataFp16 - realcutcptFp16, 'f', -1, 32)+"TF"
									backup.LeftFp32 = strconv.FormatFloat(remainfp32(cptdataFp32 - realcutcptFp32), 'f', -1, 32)+"TF"
									backup.VxpuInfo.AllocatedMemory = strconv.Itoa(int(reqmem))+"MB"
	
									// 标准单位
									backup.VxpuInfo.AllocateInt4 = strconv.FormatFloat(realcutcptInt4, 'f', -1, 32)+"TOPS"
									backup.VxpuInfo.AllocateInt8 = strconv.FormatFloat(realcutcptInt8, 'f', -1, 32)+"TOPS"
									backup.VxpuInfo.AllocateFp16 = strconv.FormatFloat(realcutcptFp16, 'f', -1, 32)+"TF"
									backup.VxpuInfo.AllocateFp32 = strconv.FormatFloat(realcutcptFp32, 'f', -1, 32)+"TF"
	
								}
	
								// memory -》precision -》flops 判断完了，缓存分配结果
								backuparray = append(backuparray, backup)
								successcnt++
	
								// 更新当前metrics信息
								tmp_metrics.Xpus[index].Vxpus = append(tmp_metrics.Xpus[index].Vxpus, backuparray[nodenum].VxpuInfo)
								tmp_metrics.Xpus[index].LeftMemory = backuparray[nodenum].LeftMemory
								// log.Println("第",nodenum,"个vxpu, tmp_metrics现在的内存是",tmp_metrics.Xpus[index].LeftMemory)
								tmp_metrics.Xpus[index].LeftInt4 = backuparray[nodenum].LeftInt4 
								tmp_metrics.Xpus[index].LeftInt8 = backuparray[nodenum].LeftInt8 
								tmp_metrics.Xpus[index].LeftFp16 = backuparray[nodenum].LeftFp16
								tmp_metrics.Xpus[index].LeftFp32 = backuparray[nodenum].LeftFp32 
								logrus.Info("第",nodenum,"个vxpu, tmp_metrics现存算力都是",tmp_metrics.Xpus[index].LeftInt4,
																					tmp_metrics.Xpus[index].LeftInt8,
																					tmp_metrics.Xpus[index].LeftFp16,
																					tmp_metrics.Xpus[index].LeftFp32)
	
								
								// 分配完了，不在循环UUID，跳出
								break
							}
							
						}
	
						// UUID循环完了都没对上
						if index == int(metrics_allocate.XpuAttitude[typenum].XpuNum) {
							allocresp.Result = fmt.Sprintf("Cannot fine Parameter in applyList[%d], XpuUUID: %s", nodenum+1,
												allocreq3.ApplyList[nodenum].Extendinfo["xpuUUID"])
							goto CONSTRUCT_RESP
						}
					}
				}
			}
		}
	
		logrus.Trace("开始检查容器")
		
		// oldcnt_max 存储目前应该被加的容器第一个编号
		oldcnt_max = int(Cnt_num(metrics_allocate, port_start))
	
		if successcnt == paralen {
			finalcnt = 0 //最终成功的容器
			for nodenum := 0; nodenum < paralen; nodenum++ {
				create_cnt_num = append(create_cnt_num, oldcnt_max + nodenum)
				index := backuparray[nodenum].Index
				// instance := uint(create_cnt_num[nodenum])
				// instanceno := fmt.Sprintf("%05d", instance)
				// instanceno := fmt.Sprintf("%05d", 30100+instance)
	
				totalinstance := int(create_cnt_num[nodenum])
				// containersequence := fmt.Sprintf("%05d", portno_start+totalinstance)
				containersequence := fmt.Sprintf("%05d", totalinstance)
	
				backuparray[nodenum].VxpuInfo.XpuUUID = metrics_allocate.Xpus[index].XpuUUID+"-"+containersequence
				backuparray[nodenum].VxpuInfo.State = "Avaliable"
				// backuparray[nodenum].VxpuInfo.ContainerName = "yxdtest_"+containersequence
				backuparray[nodenum].VxpuInfo.ContainerName = c_name + containersequence
				backuparray[nodenum].VxpuInfo.ContainerPort = containersequence
				backuparray[nodenum].VxpuInfo.ContainerIP = metrics_allocate.HnIdentifyInfo.HostIP
	
				// cnt_re, cnt_err := CreateContainer(int(create_cnt_num[nodenum]), 
				// 					metrics_allocate.Xpus[index].XpuUUID, backuparray[nodenum].Percentage, 
				// 					backuparray[nodenum].Memory)
				cnt_re, cnt_err := CreateContainer(containersequence, 
									metrics_allocate.Xpus[index].XpuUUID, 
									backuparray[nodenum].VxpuInfo.XpuUUID,
									backuparray[nodenum].Percentage, 
									backuparray[nodenum].Memory, backuparray[nodenum].VxpuInfo.ContainerName, image)
				logrus.Info("第",nodenum,"个vxpu, 容器生成结果是",cnt_re, cnt_err)
	
				if  cnt_re == "no"{
					if nodenum == 0 {
						allocresp.Result = fmt.Sprintf("Cannot create container in applyList[%d] for [%s]", nodenum, cnt_err)
						goto CONSTRUCT_RESP
					}else {
						for cancel_cnt := create_cnt_num[nodenum-1]; cancel_cnt >= oldcnt_max; cancel_cnt-- {
							DeleteContainer(backuparray[nodenum].VxpuInfo.ContainerName, tmp_file_path)
						}
						allocresp.Result = fmt.Sprintf("Cannot create container in applyList[%d] for [%s]", nodenum, cnt_err)
						goto CONSTRUCT_RESP
					}
				}
				finalcnt ++
			}
		}
	
		// 要分配的都对上了，开始启动容器
		if finalcnt == paralen {
			logrus.Trace("开始正常生成容器")
			for nodenum := 0; nodenum < paralen; nodenum++ {
				index := backuparray[nodenum].Index
	
				var vxpuinfo Allocresplist
	
				vxpuinfo.TransId = allocreq3.ApplyList[nodenum].TransId
				vxpuinfo.VxpuUUID = backuparray[nodenum].VxpuInfo.XpuUUID
	
				// extendinfo作为map可以随时更新
				vxpuinfo.ExtendInfo = make(map[string]string)
				vxpuinfo.ExtendInfo["containerName"] = backuparray[nodenum].VxpuInfo.ContainerName 
				vxpuinfo.ExtendInfo["containerPort"] = backuparray[nodenum].VxpuInfo.ContainerPort 
				// vxpuinfo.ExtendInfo["containerIP"] = backuparray[nodenum].VxpuInfo.ContainerIP
				vxpuinfo.ExtendInfo["containerIP"] = dataPathIP.addr[0]
				logrus.Info("%v: name %v, addr %v\n", backuparray[nodenum].VxpuInfo.ContainerName, dataPathIP.name[0], dataPathIP.addr[0])

				vxpuinfo.ExtendInfo["memory"] = backuparray[nodenum].VxpuInfo.AllocatedMemory
				vxpuinfo.ExtendInfo["precision"] = allocreq3.ApplyList[nodenum].Precision
				if allocreq3.ApplyList[nodenum].Precision == "Int4" {
					vxpuinfo.ExtendInfo["flops"] = trans_T2G(backuparray[nodenum].VxpuInfo.AllocateInt4)	
				}else if allocreq3.ApplyList[nodenum].Precision == "Int8" {
					vxpuinfo.ExtendInfo["flops"] = trans_T2G(backuparray[nodenum].VxpuInfo.AllocateInt8)	
				} else if allocreq3.ApplyList[nodenum].Precision == "Fp16" {
					vxpuinfo.ExtendInfo["flops"] = trans_T2G(backuparray[nodenum].VxpuInfo.AllocateFp16)	
				} else if allocreq3.ApplyList[nodenum].Precision == "Fp32" {
					vxpuinfo.ExtendInfo["flops"] = trans_T2G(backuparray[nodenum].VxpuInfo.AllocateFp32)	
				}
				vxpuinfo.ExtendInfo["type"] = "GPU"
	
				allocresp.VxpuInfo = append(allocresp.VxpuInfo, vxpuinfo)
	
				//update info
				metrics_allocate.Xpus[index].VxpuNum++
				metrics_allocate.HnIdentifyInfo.TotalVxpuNum++
	
				metrics_allocate.Xpus[index].Vxpus = append(metrics_allocate.Xpus[index].Vxpus, backuparray[nodenum].VxpuInfo)
				metrics_allocate.Xpus[index].LeftMemory = backuparray[nodenum].LeftMemory
				// log.Println("第",nodenum,"个vxpu, metrics现在的内存是",metrics_allocate.Xpus[index].LeftMemory)
				metrics_allocate.Xpus[index].LeftInt4 = backuparray[nodenum].LeftInt4 
				metrics_allocate.Xpus[index].LeftInt8 = backuparray[nodenum].LeftInt8 
				metrics_allocate.Xpus[index].LeftFp16 = backuparray[nodenum].LeftFp16
				metrics_allocate.Xpus[index].LeftFp32 = backuparray[nodenum].LeftFp32 
				logrus.Info("第",nodenum,"个vxpu, metrics现存算力都是",metrics_allocate.Xpus[index].LeftInt4,
																	metrics_allocate.Xpus[index].LeftInt8,
																	metrics_allocate.Xpus[index].LeftFp16,
																	metrics_allocate.Xpus[index].LeftFp32)
	
				// 记录当前改变的xpu的uuid
				changenum := len(changexpu)
				if changenum == 0{
					changexpu = append(changexpu, metrics_allocate.Xpus[index].XpuUUID)
				}else {
					no := 0
					for i := 0; i <  changenum; i++ {
						if metrics_allocate.Xpus[index].XpuUUID == changexpu[i] {
							break
						}else {
							no ++
						}
					}
					if no ==  changenum {
						changexpu = append(changexpu, metrics_allocate.Xpus[index].XpuUUID)
						logrus.Trace("当前changexpu有更新")
					}
				}
				
			}
			// 给返回结果
			allocresp.State = "true"
			allocresp.Result = "no error"
			
		} else {
			allocresp.Result = "Fail"
		}
	
		// 只更新当前有变动的xpu的当前状态，前面后面更新的都不管
		for i := 0;i < len(changexpu); i++ {
			for index := 0; index < len(metrics_allocate.Xpus); index ++{
				if changexpu[i] == metrics_allocate.Xpus[index].XpuUUID {
					var afterassign AfterAssignXpu
					afterassign.XpuName = metrics_allocate.Xpus[index].XpuName
					afterassign.XpuUUID = metrics_allocate.Xpus[index].XpuUUID
					afterassign.XpuDeviceID = metrics_allocate.Xpus[index].XpuDeviceID
					afterassign.State = metrics_allocate.Xpus[index].State
					afterassign.LeftMemory = metrics_allocate.Xpus[index].LeftMemory
					afterassign.LeftInt4 = metrics_allocate.Xpus[index].LeftInt4
					afterassign.LeftInt8 = metrics_allocate.Xpus[index].LeftInt8
					afterassign.LeftFp16 = metrics_allocate.Xpus[index].LeftFp16
					afterassign.VxpuNum = metrics_allocate.Xpus[index].VxpuNum
					afterassign.LeftFp32 = metrics_allocate.Xpus[index].LeftFp32
					allocresp.AfterAssign = append(allocresp.AfterAssign, afterassign)
				}
			}
		}
	
		// 不管上面goto啥了，走到这都得进来
		CONSTRUCT_RESP:
		// 返回allocresp
		json.NewEncoder(w).Encode(allocresp)
	
		// //write back to json file
		// filePtr, err := os.OpenFile(xpu_info_file, os.O_RDWR, os.ModePerm)
		// if err != nil {
		// 	log.Println("文件打开失败 [Err:%s]", err.Error())
		// 	return
		// }
		// defer filePtr.Close()

		// encoder := json.NewEncoder(fp_allocate)
		// err = encoder.Encode(metrics_allocate)

		encoder, err := json.Marshal(metrics_allocate)
		err = ioutil.WriteFile(xpu_info_file, encoder, 0777)

		if err != nil {
			logrus.Error("allocate编码错误", err.Error())
		} else {
			logrus.Trace("allocate编码成功")
		}
		
		// 解锁
		if err := syscall.Flock(int(fp_allocate.Fd()), syscall.LOCK_UN); err != nil {
			logrus.Error("allocate解锁失败", err)
		}else{
			logrus.Trace("allocate解锁成功")
		}
		// fp_allocate.Close()  // 从系统调用来看，应该是这里的问题
		defer fp_allocate.Close()
	})

	http.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("======= release ========")
		// 从文件获取当前设备信息并且给文件上锁
		var metrics_release Metrics
		// for i := 0; i < len(metrics_release.Xpus); i++ {
		// 	tmp_metrics.Xpus = append(tmp_metrics.Xpus, metrics_release.Xpus[i])
		// }
		fp_release, err := os.OpenFile(xpu_info_file, os.O_RDWR, os.ModePerm)
		if err != nil {
			logrus.Error("文件打开失败 [Err:%s]", err.Error())
			return
		}
		// 非阻塞模式下加共享锁
		if err := syscall.Flock(int(fp_release.Fd()), syscall.LOCK_EX); err != nil {
			logrus.Error("release上锁失败", err)
		}else{
			logrus.Trace("release上锁成功")
		}

		// defer fp_release.Close()
		// 创建json解码器
		decoder := json.NewDecoder(fp_release)
		err = decoder.Decode(&metrics_release)

		logrus.Trace("-----------------开始释放进程---------------------")

		var releasereq ReleaseReq
		var releaseresp ReleaseResp
		json.NewDecoder(r.Body).Decode(&releasereq)
		logrus.Trace(releasereq)
		
		// 记录当前请求下有哪些xpu改变，注意是xpu不是vxpu
		var changexpu []string
	
		// 记录成功销毁了几个容器
		listnum := 0
	
		// 记录请求销毁几个容器
		paralen := len(releasereq.ReleaseList)
		logrus.Trace(paralen)

		releaseresp.State = "false"
		releaseresp.Result = "No match vxpu"  // errormsg
		if paralen == 0 {
			releaseresp.Result = "Parameter list in empty"
			goto CONSTRUCT_RESP
		}
	
		for nodenum := 0; nodenum < paralen; nodenum++ {
			//准备返回值，每个请求都必须必须返回，不能跳过或为空
			var releaseresp_vxpuinfo Releaseresplist
			releaseresp_vxpuinfo.VxpuUUID = "" // 默认为空
			releaseresp_vxpuinfo.State = "false"  // 默认false
	
			if releasereq.ReleaseList[nodenum].VxpuUUID == "" {
				str_result := fmt.Sprintf("Parameter uncorrect in releaseList[%d]", nodenum)
				releaseresp.Result = releaseresp.Result + ";" + str_result
				continue
			}else{
				// uuid 检查完毕，不为空，把请求的uuid赋值给他
				releaseresp_vxpuinfo.VxpuUUID = releasereq.ReleaseList[nodenum].VxpuUUID
				// 正式检查开始
				for typenum := 0; typenum < metrics_release.HnIdentifyInfo.XpuTypeNum; typenum++ {
					// 循环xpu类型
					index := 0 //index提前声明是为了下面循环满了不报错undefined
					for ; index < int(metrics_release.XpuAttitude[typenum].XpuNum); index++ {
						logrus.Trace("第",nodenum,"个vxpu，此时index=",index)
						// 循环xpu类型下的每一个xpu
	
						// 获得要释放的容器的id
						vxpuuuid := releasereq.ReleaseList[nodenum].VxpuUUID
	
						// 获得要释放的xpu的id  FIXME:
						xpuuuid := vxpuuuid[:len(vxpuuuid)-6]
	
						// 目前的内存资源
						totalint4cpt := metrics_release.XpuAttitude[typenum].TotalInt4 // 
						totalint4cptdata,_ := strconv.ParseFloat(totalint4cpt[:len(totalint4cpt)-4], 32)
						totalint8cpt := metrics_release.XpuAttitude[typenum].TotalInt8 // 
						totalint8cptdata,_ := strconv.ParseFloat(totalint8cpt[:len(totalint8cpt)-4], 32)
						totalfp16cpt := metrics_release.XpuAttitude[typenum].TotalFP16 // 
						totalfp16cptdata,_ := strconv.ParseFloat(totalfp16cpt[:len(totalfp16cpt)-2], 32)
						totalfp32cpt := metrics_release.XpuAttitude[typenum].TotalFP32 // 
						totalfp32cptdata,_ := strconv.ParseFloat(totalfp32cpt[:len(totalfp32cpt)-2], 32)
						totalfp32cptdata = remainfp32(totalfp32cptdata)
	
						if xpuuuid == metrics_release.Xpus[index].XpuUUID {
							// 匹配上xpu的uuid
							// 循环"Vxpus"
							for instance := 0; instance < int(metrics_release.Xpus[index].VxpuNum); instance++ {
								// 循环vxpu
								logrus.Trace("yes5:%s,%s",releasereq.ReleaseList[nodenum].VxpuUUID,metrics_release.Xpus[index].Vxpus[instance].XpuUUID)
								if releasereq.ReleaseList[nodenum].VxpuUUID == metrics_release.Xpus[index].Vxpus[instance].XpuUUID {
									//对比含有末尾四位数字的vxpuuuid
									//匹配上了，开始分配
	
									// 检查内存
									leftmem := metrics_release.Xpus[index].LeftMemory
									leftmemdata,_ := strconv.Atoi(leftmem[:len(leftmem)-2]) // 去掉单位
	
									allocmem := metrics_release.Xpus[index].Vxpus[instance].AllocatedMemory 
									allocmemdata,_ := strconv.Atoi(allocmem[:len(allocmem)-2])
	
									totalmem := metrics_release.XpuAttitude[typenum].TotalMemory
									totalmemdata,_ := strconv.Atoi(totalmem[:len(totalmem)-2])
	
									logrus.Trace("第",nodenum,"个vxpu，此时leftmemdata=",leftmemdata)
									logrus.Trace("第",nodenum,"个vxpu，此时allocmemdata=",allocmemdata)
									logrus.Trace("第",nodenum,"个vxpu，此时totalmemdata=",totalmemdata)
	
									if (leftmemdata + allocmemdata) > totalmemdata {
										str_result := fmt.Sprintf("Memory over allocate in releaseList[%d]", nodenum)
										releaseresp.Result = releaseresp.Result + ";" + str_result
										continue
									} else {
	
										leftcptInt4 := metrics_release.Xpus[index].LeftInt4 // 
										cptdataInt4,_ := strconv.ParseFloat(leftcptInt4[:len(leftcptInt4)-4], 32) // Int4
										leftcptInt8 := metrics_release.Xpus[index].LeftInt8
										cptdataInt8,_ := strconv.ParseFloat(leftcptInt8[:len(leftcptInt8)-4], 32) // Int8
										leftcptFp16 := metrics_release.Xpus[index].LeftFp16
										cptdataFp16,_ := strconv.ParseFloat(leftcptFp16[:len(leftcptFp16)-2], 32) // Fp16
										leftcptFp32 := metrics_release.Xpus[index].LeftFp32
										cptdataFp32,_ := strconv.ParseFloat(leftcptFp32[:len(leftcptFp32)-2], 32) // Fp32
										cptdataFp32 = remainfp32(cptdataFp32)
										
										alloccptInt4 := metrics_release.Xpus[index].Vxpus[instance].AllocateInt4 
										alloccptdataInt4,_ := strconv.ParseFloat(alloccptInt4[:len(alloccptInt4)-4], 32)
										alloccptInt8 := metrics_release.Xpus[index].Vxpus[instance].AllocateInt8 
										alloccptdataInt8,_ := strconv.ParseFloat(alloccptInt8[:len(alloccptInt8)-4], 32)
										alloccptFp16 := metrics_release.Xpus[index].Vxpus[instance].AllocateFp16
										alloccptdataFp16,_ := strconv.ParseFloat(alloccptFp16[:len(alloccptFp16)-2], 32)
										alloccptFp32 := metrics_release.Xpus[index].Vxpus[instance].AllocateFp32
										alloccptdataFp32,_ := strconv.ParseFloat(alloccptFp32[:len(alloccptFp32)-2], 32)
										alloccptdataFp32 = remainfp32(alloccptdataFp32)
	
	
										logrus.Info("第",nodenum,"个vxpu，此时要释放的算力都是", alloccptdataInt4, alloccptdataInt8, alloccptdataFp16, alloccptdataFp32)
										
										if metrics_release.Xpus[index].Vxpus[instance].AllocateInt4 == "" &&
										metrics_release.Xpus[index].Vxpus[instance].AllocateInt8 == "" &&
										metrics_release.Xpus[index].Vxpus[instance].AllocateFp16 == "" &&
										metrics_release.Xpus[index].Vxpus[instance].AllocateFp32 == "" {
											str_result := fmt.Sprintf("Compute unsupport in releaseList[%d]", nodenum)
											releaseresp.Result = releaseresp.Result + ";" + str_result
											continue
										}
	
										if (cptdataInt4 + alloccptdataInt4) > totalint4cptdata ||
												(cptdataInt8 + alloccptdataInt8) > totalint8cptdata ||
												(cptdataFp16 + alloccptdataFp16) > totalfp16cptdata ||
												remainfp32((cptdataFp32 + alloccptdataFp32)) > totalfp32cptdata {
													str_result := fmt.Sprintf("Compute over allocate in releaseList[%d]", nodenum)
													releaseresp.Result = releaseresp.Result + ";" + str_result
													logrus.Info("第",nodenum,"个vxpu, 超出的和是",cptdataFp32 + alloccptdataFp32)
													continue
										}
	
										//delete container
										DeleteContainer(metrics_release.Xpus[index].Vxpus[instance].ContainerName, tmp_file_path)
	
										listnum++
	
										//update info
										releaseresp_vxpuinfo.State = "true" 
										releaseresp.VxpuInfo = append(releaseresp.VxpuInfo, releaseresp_vxpuinfo)
										logrus.Trace("第",nodenum,"个vxpu, 返回格式是",releaseresp.VxpuInfo)
	
										metrics_release.Xpus[index].LeftInt4 = strconv.FormatFloat(cptdataInt4 + alloccptdataInt4, 'f', -1, 32)+"TOPS"
										metrics_release.Xpus[index].LeftInt8 = strconv.FormatFloat(cptdataInt8 + alloccptdataInt8, 'f', -1, 32)+"TOPS"
										metrics_release.Xpus[index].LeftFp16 = strconv.FormatFloat(cptdataFp16 + alloccptdataFp16, 'f', -1, 32)+"TF"
										metrics_release.Xpus[index].LeftFp32 = strconv.FormatFloat(remainfp32(cptdataFp32 + alloccptdataFp32), 'f', -1, 32)+"TF"
										logrus.Info("第",nodenum,"个vxpu, metrics现存算力都是",metrics_release.Xpus[index].LeftInt4,
																						metrics_release.Xpus[index].LeftInt8,
																						metrics_release.Xpus[index].LeftFp16,
																						metrics_release.Xpus[index].LeftFp32)
										metrics_release.Xpus[index].VxpuNum = metrics_release.Xpus[index].VxpuNum - 1
										metrics_release.HnIdentifyInfo.TotalVxpuNum  = metrics_release.HnIdentifyInfo.TotalVxpuNum - 1
										metrics_release.Xpus[index].LeftMemory = strconv.Itoa(int(leftmemdata + allocmemdata))+"MB"
										logrus.Trace("第",nodenum,"个vxpu, metrics现在的内存是",metrics_release.Xpus[index].LeftMemory)
	
										metrics_release.Xpus[index].Vxpus = append(metrics_release.Xpus[index].Vxpus[:instance],metrics_release.Xpus[index].Vxpus[instance+1:]...)
									}
	
									// 记录当前改变的xpu的uuid
									changenum := len(changexpu)
									if changenum == 0{
										changexpu = append(changexpu, metrics_release.Xpus[index].XpuUUID)
										logrus.Trace("当前changexpu有更新")
									}else {
										no := 0
										for i := 0; i <  changenum; i++ {
											if metrics_release.Xpus[index].XpuUUID == changexpu[i] {
												break
											}else {
												no ++
											}
										}
										if no ==  changenum {
											changexpu = append(changexpu, metrics_release.Xpus[index].XpuUUID)
											logrus.Trace("当前changexpu有更新")
										}
									}
	
									break
								}
							}
						}
					}
					// UUID循环完了都没对上
					if index == int(metrics_release.XpuAttitude[typenum].XpuNum) &&
						releaseresp_vxpuinfo.State == "false"  {
						str_result := fmt.Sprintf("Cannot fine Parameter in applyList[%d], XpuUUID: %s", nodenum,
																			releasereq.ReleaseList[nodenum].VxpuUUID)
						releaseresp.Result = releaseresp.Result + ";" + str_result
						releaseresp.VxpuInfo = append(releaseresp.VxpuInfo, releaseresp_vxpuinfo)
						logrus.Trace("第",nodenum,"个vxpu, 返回格式是",releaseresp.VxpuInfo)
					}
				}
			}
			
		}
	
		if listnum == len(releasereq.ReleaseList) {
			releaseresp.State = "true"
			releaseresp.Result = "no error"
		}
	
		// 只更新当前有变动的xpu的当前状态，前面后面更新的都不管
		for i := 0;i < len(changexpu); i++ {
			for index := 0; index < len(metrics_release.Xpus); index ++{
				if changexpu[i] == metrics_release.Xpus[index].XpuUUID {
					var afterassign AfterAssignXpu
					afterassign.XpuName = metrics_release.Xpus[index].XpuName
					afterassign.XpuUUID = metrics_release.Xpus[index].XpuUUID
					afterassign.State = metrics_release.Xpus[index].State
					afterassign.LeftMemory = metrics_release.Xpus[index].LeftMemory
					afterassign.LeftInt4 = metrics_release.Xpus[index].LeftInt4
					afterassign.LeftInt8 = metrics_release.Xpus[index].LeftInt8
					afterassign.LeftFp16 = metrics_release.Xpus[index].LeftFp16
					afterassign.VxpuNum = metrics_release.Xpus[index].VxpuNum
					afterassign.LeftFp32 = metrics_release.Xpus[index].LeftFp32
					releaseresp.AfterAssign = append(releaseresp.AfterAssign, afterassign)
				}
			}
		}
	
	
		CONSTRUCT_RESP:
		
		json.NewEncoder(w).Encode(releaseresp)
	
		// //write back to json file
		// fp_release, err := os.OpenFile(xpu_info_file, os.O_RDWR, os.ModePerm)
		// if err != nil {
		// 	log.Println("文件打开失败 [Err:%s]", err.Error())
		// 	return
		// }
		// defer fp_release.Close()
		// 创建json解码器
		// encoder := json.NewEncoder(fp_release)
		// err = encoder.Encode(metrics_release)

		encoder, err := json.Marshal(metrics_release)
		err = ioutil.WriteFile(xpu_info_file, encoder, 0777)

		if err != nil {
			logrus.Error("release编码错误", err.Error())
		} else {
			logrus.Trace("release编码成功")
		}

		// 解锁
		if err := syscall.Flock(int(fp_release.Fd()), syscall.LOCK_UN); err != nil {
			logrus.Error("release解锁失败", err)
		}else{
			logrus.Trace("release解锁成功")
		}
		fp_release.Close()
	
	})

	var host_and_port string = hostIP+":"+port
	logrus.Info("ListenAndServer metrics: "+host_and_port+"/mettics")
	logrus.Info("ListenAndServer release: "+host_and_port+"/release")
	logrus.Info("ListenAndServer allocate: "+host_and_port+"/allocate")
	http.ListenAndServe(":"+port+"", nil) //端口设置，11.160.41.178:9406/xxxxx
}

func cpt_percent_t4(minisize float64, want float64, total float64) (float64) {
	// var minisize float64 = 1.85
	var percent float64
	eps := 0.000001
	p := (want / total) * 100 // 算力百分比
	decimal := p/minisize - math.Floor(p/minisize) // 测试是否能恰好分完
	if decimal < eps{
		percent = (math.Floor(p/minisize)) * minisize // 54*1.85 = 99.9
	}else{
		percent = (math.Floor(p/minisize) + 1) * minisize // 55*1.85 = 101.75
	}
	
	percent_sup := (math.Floor(100/minisize) + 1) * minisize  // 往上取一个浮动上限
	logrus.Trace("sup=", percent_sup)
	if (percent > 100) && (percent <= percent_sup) {
		percent = 100.0
	}
	logrus.Trace("percent=", percent)
	return percent
}

func CreateContainer(port string, uuid string, vxpu_uuid string, compute string, memory string, name string, image string) (Result string, Err string){
	Result = "yes"
	Err = "no error"
	p := compute[:len(compute)-1] // 要切割的算力
	// m := memory[:len(memory)-2] // 要切割的显存。mps不能切割显存，暂时不用

	docker_base_args := " -tid --ulimit memlock=-1:-1 --net=host "
	nvidia_args := " --runtime=nvidia " +
		" --env CUDA_VISIBLE_DEVICES=" + uuid +
		" -e NVIDIA_DRIVER_CAPABILITIES=video,compute,utility " +
		" --cap-add=SYS_PTRACE " +
		" --security-opt seccomp=unconfined " +
		" --privileged "

	vodla_server_args := " --device=/dev/infiniband/ " + 
		" -e C_PORT=" + port + 
		" -p " + port + ":" +  port + 
		" --name " + name +
		" --rm " + image

	sinian_metric_args := " -e VXPU_UUID=" + vxpu_uuid
	mps_args := " --env CUDA_MPS_ACTIVE_THREAD_PERCENTAGE=" + "\"" + p + "\"" +
		" -v /tmp/nvidia-mps:/tmp/nvidia-mps " +
		" --ipc=host "

	docker_run_cmd := "docker run " + 
		docker_base_args + 
		nvidia_args +
		vodla_server_args +
		sinian_metric_args +
		mps_args	
	
	logrus.Info(docker_run_cmd)

	err, out, errout := Shellout(docker_run_cmd)
	if err != nil {
		logrus.Error("出错：", out)
		logrus.Error("error: %v\n", err)
		Result = "no"
		Err = errout
	}

	return Result, Err
}

func DeleteContainer(containername string, sharepathfather string) {

	cmd := "docker rm -f "+containername+" &"
	logrus.Trace(cmd)
	err, _, _ := ShellStart(cmd)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	//log.Println("--- stdout ---")
	//log.Println(out)
	//log.Println("--- stderr ---")
	//log.Println(errout)
	// sharepathfather := *arg_share_path
    // sharepath := "/home1/yixiaodie.yxd/code/"+containername
	sharepath := sharepathfather+containername
	// sharepath := "/tmp/vxpu/"+containername
    cmd0 := "rm -rf "+sharepath+" &"
	logrus.Trace(cmd0)
    Shellout(cmd0)
}
const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd := exec.Command(ShellToUse, "-c", command)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    return err, stdout.String(), stderr.String()
}

func ShellStart(command string) (error, string, string) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd := exec.Command(ShellToUse, "-c", command)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Start()
    return err, stdout.String(), stderr.String()
}

// 容器编号规则：接着现存容器最大的编号往后加一
func Cnt_num(metrics Metrics, start int64) (int64) {
	logrus.Trace("计算容器编号")
	var sum int64 = 0
	
	for i := 0; i < metrics.HnIdentifyInfo.XpuTypeNum; i++ {
		for k := 0; k < int(metrics.XpuAttitude[i].XpuNum); k++{
			var a int64 = 0
			for j := 0; j < int(metrics.Xpus[k].VxpuNum); j++ {
				b := metrics.Xpus[k].Vxpus[j].XpuUUID
				logrus.Trace("vxpuuuid=", b)
	
				c := b[len(b)-5:]
				logrus.Trace("vxpuuuid后5位=", c)
	
				d,_ := strconv.ParseInt(c, 10, 0)
				logrus.Trace("vxpu整数序列=", d)
	
				if d > a {
					a = d
				}
				logrus.Trace("a=", a)
			}
			if a > sum {
				sum = a
			}
			logrus.Trace("sum=", sum)
		}
	}
	if metrics.HnIdentifyInfo.TotalVxpuNum == 0 {
		sum = start
	}else{
		sum ++
	}
	logrus.Trace("现在该加第%d个容器", sum)
	return sum
}

// 保留float64
func remainfp32(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value 
}

// 把TOPS和TF都转为G
func trans_T2G(origin_T string) (string) {

	end4 := origin_T[len(origin_T)-4:]
	end2 := origin_T[len(origin_T)-2:]
	var re string
	if end4 == "TOPS" {
		origin_T = origin_T[:len(origin_T)-4]
		count,_ := strconv.ParseFloat(origin_T, 32)
		count = count * 1000
		re = strconv.FormatFloat(count, 'f', -1, 32)+"G"
	}else if end2 == "TF" {
		origin_T = origin_T[:len(origin_T)-2]
		count,_ := strconv.ParseFloat(origin_T, 32)
		count = remainfp32(count)
		count = count * 1000
		re = strconv.FormatFloat(count, 'f', -1, 32)+"G"
	}else {
		logrus.Error("error: unit failed")
	}
	return re
}

// 根据gpu不同型号进行不同算力的切割
func cut_minisize(xpuname string) (float64) {
	var re float64 = 0
	XputypeInfo := make(map[string]float64)

	// 添加不同类型的算力最小粒度
	XputypeInfo["Tesla T4"] = 5
	XputypeInfo["Tesla V100"] = 2.5
	// XputypeInfo["Tesla A100"] = 1.85
	XputypeInfo["NVIDIA A100-SXM4-40GB"] = 1.85

	// 匹配
	for key, value := range XputypeInfo {
		if xpuname == key {
			re = value
			break
		}
	}
	return re
}

func PathExists(path string) (bool, error) {

    _, err := os.Stat(path, )
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

type IpInfo struct {
	name [2]string
	addr [2]string
}

func localAddresses(ips *IpInfo) error {
	ip_cnt := 0
    ifaces, err := net.Interfaces()
    if err != nil {
        logrus.Info(fmt.Errorf("localAddresses: %v\n", err.Error()))
        return err
    }
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            logrus.Info(fmt.Errorf("localAddresses: %v\n", err.Error()))
            continue
        }

        for _, a := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if (strings.Compare(i.Name,  "bond0") == 0) ||
				   (strings.Compare(i.Name, "ens5f1") == 0) {
					   if ip_cnt < 2 {
					   	ips.name[ip_cnt] = i.Name
					   	ips.addr[ip_cnt] = ipnet.IP.String()
					   	ip_cnt++
				   	   }
				}
			}
		}
        }
    }
    return err
}

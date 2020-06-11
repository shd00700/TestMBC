package mbclient


import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Scrclr(){
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

}
const (
	Init        = "Init"
	ModbusError = "ModbusError"
	Ok          = "Ok"
	Disconnect  = "Disconnect"
)

//MBClient config
type MBClient struct {
	IP      string
	Port    int
	Timeout time.Duration
	Conn    net.Conn
}


//state show for error

// NewClient creates a new Modbus Client config.
func NewClient(IP string, port int, timeout time.Duration) *MBClient {
	print("클라이언트")
	m := &MBClient{}
	m.IP = IP
	m.Port = port
	m.Timeout = timeout

	return m
}

//Open modbus tcp connetion
func (m *MBClient) Open() error {
	print("오픈")
	addr := m.IP + ":" + strconv.Itoa(m.Port)
	// var err error
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		print("하하")
		log.Panicln(err)
	}
	m.Conn = conn

	return nil
}

//Close modbus tcp connetion
func (m *MBClient) Close() {
	print("ㅎㅎ")
	//os.Exit(12)
	//print("gg")
	if m.Conn != nil {
		print("gg")
		m.Conn.Close()
	}
}

//IsConnected for check modbus connetection
func (m *MBClient) IsConnected() bool {
	if m.Conn != nil {
		return true
	}
	return false
}

//Qurry make a modbus tcp qurry
func Qurry(conn net.Conn,timeout time.Duration, pdu []byte) ([]byte, error) {
	if conn == nil {
		return []byte{}, fmt.Errorf(Disconnect)
	}
	header := []byte{0, 0, 0, 0, byte(len(pdu) << 10), byte(len(pdu))}
	wbuf := append(header, pdu...)
	//write
	_, err := conn.Write([]byte(wbuf))
	if err != nil {
		return nil, fmt.Errorf(Disconnect)
	}

	//read
	rbuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))
	len, err := conn.Read(rbuf)
	if err != nil {
		return nil, fmt.Errorf(Disconnect)
	}
	if len < 10 {
		return nil, fmt.Errorf(ModbusError)
	}
	return rbuf[6:len], nil
}

//ReadCoil mdbus function 1 qurry and return []uint16
func (m *MBClient) ReadCoil(id uint8, addr uint16, leng uint16) ([]int, error) {
	pdu := []byte{id, 0x01, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []int{}, err
	}
	//convert
	result := []int{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, 1)
			} else {
				result = append(result, 0)
			}
		}
	}
	result = result[:leng]
	return result, nil
}

//ReadCoilIn mdbus function 2 qurry and return []uint16
func (m *MBClient) ReadCoilIn(id uint8, addr uint16, leng uint16) ([]int, error) {

	pdu := []byte{id, 0x02, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []int{}, err
	}

	//convert
	result := []int{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, 1)
			} else {
				result = append(result, 0)
			}
		}
	}
	result = result[:leng]

	return result, nil
}

//ReadReg mdbus function 3 qurry and return []uint16
func (m *MBClient) ReadReg(id uint8, addr uint16, leng uint16) ([]uint16, error) {

	pdu := []byte{id, 0x03, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []uint16{}, err
	}
	//convert
	result := []uint16{}
	for i := 0; i < int(leng); i++ {
		var b uint16
		b = uint16(res[i*2+3]) << 8
		b |= uint16(res[i*2+4])
		result = append(result, b)
	}

	return result, nil
}

//ReadRegIn mdbus function 4 qurry and return []uint16
func (m *MBClient) ReadRegIn(id uint8, addr uint16, leng uint16) ([]uint16, error) {

	pdu := []byte{id, 0x04, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []uint16{}, err
	}

	//convert
	result := []uint16{}
	for i := 0; i < int(leng); i++ {
		var b uint16
		b = uint16(res[i*2+3]) << 8
		b |= uint16(res[i*2+4])
		result = append(result, b)
	}

	return result, nil
}

//WriteCoil mdbus function 5 qurry and return []uint16
func (m *MBClient) WriteCoil(id uint8, addr uint16, data bool) error {

	var pdu = []byte{}
	if data == true {
		pdu = []byte{id, 0x5, byte(addr >> 8), byte(addr), 0xff, 0x00}
	} else {
		pdu = []byte{id, 0x5, byte(addr >> 8), byte(addr), 0x00, 0x00}
	}

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return err
	}

	return nil
}

//WriteReg mdbus function 6 qurry and return []uint16
func (m *MBClient) WriteReg(id uint8, addr uint16, data uint16) error {

	pdu := []byte{id, 0x06, byte(addr >> 8), byte(addr), byte(data >> 8), byte(data)}

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return err
	}

	return nil
}

//WriteCoils mdbus function 15(0x0f) qurry and return []uint16
func (m *MBClient) WriteCoils(id uint8, addr uint16, data []string) error {
	pdu := []byte{}
	if len(data)%8 == 0 {
		pdu = []byte{id, 0x0f, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data) / 8)}
	} else {
		pdu = []byte{id, 0x0f, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data)/8) + 1}
	}
	var tbuf byte
	for i := 0; i < len(data); i++ {
		pb, _ := strconv.ParseBool(data[i])
		pa,_ := strconv.Atoi(data[i])
		fmt.Println("alias  ", addr, ": ", pa)
		if pb {
			tbuf |= byte(1 << uint(i%8))
		}

		if (i+1)%8 == 0 || i == len(data)-1 {
			pdu = append(pdu, tbuf)
			tbuf = 0
		}
		addr++
	}
	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return err
	}
	return nil
}
//WriteRegs mdbus function 16(0x10) qurry and return []uint16


func (m *MBClient) WriteRegs(id uint8, addr uint16, data []string)  error {


	//var data []byte
	pdu := []byte{id, 0x10, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data)) * 2}
	for i := 0; i < len(data); i++ {
		pi, _ := strconv.ParseUint(data[i], 10, 16)
		fmt.Println("alias  ",addr, ": ", pi)
		pdu = append(pdu, byte(pi>>8))
		pdu = append(pdu, byte(pi))
		//fmt.Println(pdu)
		addr++
	}


	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil


		}
		return err
	}
	return nil
}

func Error(){
	println("-------------------------------------")
	fmt.Println("error")
	println("-------------------------------------")
	fmt.Println("You entered it iolncorrectly.\nReturn to the Output Coils menu.")
	println("-------------------------------------\n\n\n")
	fmt.Print("1 : Back\n","2 : Main menu\n\n\n")
	fmt.Print("Select number Enter:")
}
func Continue(){
	println("\n\n\n-------------------------------------")
	fmt.Println("1 : Back","\n2 : Main menu")
	println("-------------------------------------\n\n\n")
	fmt.Print("Select number Enter:")
}

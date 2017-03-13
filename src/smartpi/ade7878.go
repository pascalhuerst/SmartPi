/*
    Copyright (C) Jens Ramhorst
	  This file is part of SmartPi.
    SmartPi is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    SmartPi is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with SmartPi.  If not, see <http://www.gnu.org/licenses/>.
    Diese Datei ist Teil von SmartPi.
    SmartPi ist Freie Software: Sie können es unter den Bedingungen
    der GNU General Public License, wie von der Free Software Foundation,
    Version 3 der Lizenz oder (nach Ihrer Wahl) jeder späteren
    veröffentlichten Version, weiterverbreiten und/oder modifizieren.
    SmartPi wird in der Hoffnung, dass es nützlich sein wird, aber
    OHNE JEDE GEWÄHRLEISTUNG, bereitgestellt; sogar ohne die implizite
    Gewährleistung der MARKTFÄHIGKEIT oder EIGNUNG FÜR EINEN BESTIMMTEN ZWECK.
    Siehe die GNU General Public License für weitere Details.
    Sie sollten eine Kopie der GNU General Public License zusammen mit diesem
    Programm erhalten haben. Wenn nicht, siehe <http://www.gnu.org/licenses/>.
*/

package smartpi

import (
	"golang.org/x/exp/io/i2c"
  "github.com/nathan-osman/go-rpigpio"
  "time"
  "math"
	"fmt"
)


const (
  I2C_DEVICE = "/dev/i2c-1"
  ADE7878_ADDR = 0x38
  SAMPLES = 100
  ADE7878_CLOCK float32 = 256000
	FACTOR_CIRCLE float32 = 360
	VAL float32 = math.Pi / 180.0
  FACTOR_1 int = 256;
  FACTOR_2 int = 65536;
  FACTOR_3 int = 16777216;
  RMS_FACTOR_VOLTAGE float32 = 2427873
  CURRENT_RESISTOR_A float32 = 7.07107
  CURRENT_RESISTOR_B float32 = 7.07107
  CURRENT_RESISTOR_C float32 = 7.07107
  CURRENT_RESISTOR_N float32 = 7.07107
  CURRENT_CLAMP_FACTOR_A float32 = 0.05
  CURRENT_CLAMP_FACTOR_B float32 = 0.05
  CURRENT_CLAMP_FACTOR_C float32 = 0.05
  CURRENT_CLAMP_FACTOR_N float32 = 0.05
  OFFSET_CURRENT_A float32 = 0.97129167
  OFFSET_CURRENT_B float32 = 0.97129167
  OFFSET_CURRENT_C float32 = 0.97129167
  OFFSET_CURRENT_N float32 = 0.97129167
  OFFSET_VOLTAGE_A float32 = 1.0
  OFFSET_VOLTAGE_B float32 = 1.0
  OFFSET_VOLTAGE_C float32 = 1.0
	POWER_CORRECTION_FACTOR_A float32 = 0.019413
	POWER_CORRECTION_FACTOR_B float32 = 0.019413
	POWER_CORRECTION_FACTOR_C float32 = 0.019413
)


var (
	rms_factor_current float32
)


func resetADE7878() {
  println("RESET")
  p, err := rpi.OpenPin(4, rpi.OUT)
  if err != nil {
      panic(err)
  }
  defer p.Close()
  p.Write(rpi.LOW)
  time.Sleep(time.Second)
  p.Write(rpi.HIGH)
}

func initPiForADE7878() {
/*
  p, err := rpi.OpenPin(4, rpi.OUT)
  if err != nil {
      panic(err)
  }
  defer p.Close()
  p.Write(rpi.HIGH)*/
}


func InitADE7878(c *Config) (*i2c.Device, error)  {

  var dataAddress []byte
  dataAddress = make([]byte, 3)
  var i2cLock []byte
  i2cLock = make([]byte, 1)


	d, err := i2c.Open(&i2c.Devfs{Dev: I2C_DEVICE}, ADE7878_ADDR)
  if err != nil {
      panic(err)
  }

	dataAddress[0] = 0xEC;//0xEC01 (CONFIG2-REGISTER)
	dataAddress[1] = 0x01;
	dataAddress[2] = 0x02;//00000010 --> Bedeutet I2C-Lock (I2C ist nun die gewählte Übertragungsart)

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

  dataAddress = make([]byte, 1)
  dataAddress[0] = 0xEC;//0xEC01 (CONFIG2-REGISTER)

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

  err = d.Read(i2cLock)
  if err != nil {
      panic(err)
  }



	dataAddress = make([]byte, 3)
  dataAddress[0] = 0xE7;//0xE7FE writeprotection
	dataAddress[1] = 0xFE;
	dataAddress[2] = 0xAD;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }
	dataAddress[0] = 0xE7;//0xE7E3 writeprotection OFF
	dataAddress[1] = 0xE3;
	dataAddress[2] = 0x00;

	err = d.Write(dataAddress)
	if err != nil {
			panic(err)
	}

	// dataAddress = make([]byte, 6)
  // dataAddress[0] = 0x43;//0x43B6 (HPFDIS-REGISTER)
	// dataAddress[1] = 0xB6;
	// dataAddress[2] = 0x00;
	// dataAddress[3] = 0x00;
	// dataAddress[4] = 0x00;
	// dataAddress[5] = 0x00;
	//
	// err = d.Write(dataAddress)
  // if err != nil {
  //     panic(err)
  // }


	// set the right power frequency to the COMPMODE-REGISTER
	dataAddress = make([]byte, 4)
	dataAddress[0] = 0xE6;//0xE60E (COMPMODE-REGISTER)
	dataAddress[1] = 0x0E;
	if c.Powerfrequency == 60 {
		dataAddress[2] = 0x41;
		dataAddress[3] = 0xFF;
	} else {
		dataAddress[2] = 0x01;
		dataAddress[3] = 0xFF;
	}
	err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }



	dataAddress = make([]byte, 5)
  dataAddress[0] = 0x43;//0x43B5 (DICOEFF-REGISTER)
	dataAddress[1] = 0xB5;
	dataAddress[2] = 0xFF;
	dataAddress[3] = 0x80;
	dataAddress[4] = 0x00;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x43AB (WTHR1-REGISTER)
	dataAddress[1] = 0xAB;
	dataAddress[2] = 0x00;
	dataAddress[3] = 0x00;
	dataAddress[4] = 0x00;
	dataAddress[5] = 0x17;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x43AC (WTHR0-REGISTER)
	dataAddress[1] = 0xAC;
	dataAddress[2] = 0x00;
	dataAddress[3] = 0x85;
	dataAddress[4] = 0x60;
	dataAddress[5] = 0x16;


  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }
	//
	// dataAddress = make([]byte, 6)
  // dataAddress[0] = 0x43;//0x43AD (VARTHR1-REGISTER)
	// dataAddress[1] = 0xAD;
	// dataAddress[2] = 0x17;
	// dataAddress[3] = 0x85;
	// dataAddress[4] = 0x60;
	// dataAddress[5] = 0x16;
	//
	//
  // err = d.Write(dataAddress)
  // if err != nil {
  //     panic(err)
  // }
	//
	// dataAddress = make([]byte, 6)
  // dataAddress[0] = 0x43;//0x43AE (VARTHR0-REGISTER)
	// dataAddress[1] = 0xAE;
	// dataAddress[2] = 0x17;
	// dataAddress[3] = 0x85;
	// dataAddress[4] = 0x60;
	// dataAddress[5] = 0x16;
	//
	//
  // err = d.Write(dataAddress)
  // if err != nil {
  //     panic(err)
  // }
	//
	// dataAddress = make([]byte, 6)
  // dataAddress[0] = 0x43;//0x43A9 (VATHR1-REGISTER)
	// dataAddress[1] = 0xA9;
	// dataAddress[2] = 0x17;
	// dataAddress[3] = 0x85;
	// dataAddress[4] = 0x60;
	// dataAddress[5] = 0x16;
	//
	//
  // err = d.Write(dataAddress)
  // if err != nil {
  //     panic(err)
  // }
	//
	// dataAddress = make([]byte, 6)
  // dataAddress[0] = 0x43;//0x43AA (VARTHR0-REGISTER)
	// dataAddress[1] = 0xAA;
	// dataAddress[2] = 0x17;
	// dataAddress[3] = 0x85;
	// dataAddress[4] = 0x60;
	// dataAddress[5] = 0x16;
	//
	//
  // err = d.Write(dataAddress)
  // if err != nil {
  //     panic(err)
  // }

	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x43B3 (VLEVEL-REGISTER)
	dataAddress[1] = 0xB3;
	dataAddress[2] = 0x00;
	dataAddress[3] = 0x0C;
	dataAddress[4] = 0xEC;
	dataAddress[5] = 0x85;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

	time.Sleep(875 * time.Millisecond)

	// dataAddress = make([]byte, 2)
  // dataAddress[0] = 0x43;//0x4381 (AVGAIN-REGISTER)
	// dataAddress[1] = 0x81;
	// data = make([]byte, 4)
	//
	// err = d.Write(dataAddress)
	// if err != nil {
	// 		panic(err)
	// }
	// err = d.Read(data)
	// if err != nil {
	// 		panic(err)
	// }
	//
	// outcome = float32(FACTOR_3*int(data[0])+FACTOR_2*int(data[1])+FACTOR_1*int(data[2])+int(data[3]))
	// fmt.Printf("AVGAIN-REGISTER VORHER%g   %x %x %x %x \n", outcome, data[0], data[1], data[2], data[3])


	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x4381 (AVGAIN-REGISTER)
	dataAddress[1] = 0x81;
	dataAddress[2] = 0xFF;
	dataAddress[3] = 0xFC;
	dataAddress[4] = 0x1C;
	dataAddress[5] = 0xC2;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }


	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x4383 (BVGAIN-REGISTER)
	dataAddress[1] = 0x83;
	dataAddress[2] = 0xFF;
	dataAddress[3] = 0xFB;
	dataAddress[4] = 0xCA;
	dataAddress[5] = 0x60;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

	dataAddress = make([]byte, 6)
  dataAddress[0] = 0x43;//0x4385 (CVGAIN-REGISTER)
	dataAddress[1] = 0x85;
	dataAddress[2] = 0xFF;
	dataAddress[3] = 0xFC;
	dataAddress[4] = 0x12;
	dataAddress[5] = 0xDE;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }



	// Line cycle mode
	dataAddress = make([]byte, 3)
  dataAddress[0] = 0xE7;//0xE702 LCYCMODE
	dataAddress[1] = 0x02;
	dataAddress[2] = 0x0F;

	err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }

	// Line cycle mode count
	dataAddress = make([]byte, 3)
  dataAddress[0] = 0xE6;//0xE60C LINECYC
	dataAddress[1] = 0x0C;
	dataAddress[2] = 0xC8;

	err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }




	dataAddress = make([]byte, 3)
  dataAddress[0] = 0xE7;//0xE7FE writeprotection
	dataAddress[1] = 0xFE;
	dataAddress[2] = 0xAD;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }
	dataAddress[0] = 0xE7;//0xE7E3 writeprotection
	dataAddress[1] = 0xE3;
	dataAddress[2] = 0x80;

	err = d.Write(dataAddress)
	if err != nil {
			panic(err)
	}


	dataAddress = make([]byte, 4)
  dataAddress[0] = 0xE2;//0xE228 (RUN-Register)
	dataAddress[1] = 0x28;
	dataAddress[2] = 0x00;
	dataAddress[3] = 0x01;

  err = d.Write(dataAddress)
  if err != nil {
      panic(err)
  }
	return d, nil
}

func ReadoutValues(d *i2c.Device, c *Config) [25]float32 {

  var dataAddress []byte
  var data []byte
  var values [25]float32
  var outcome float32
	var err error

  initPiForADE7878()
  //resetADE7878()


	if c.Powerfrequency == 60 {
		rms_factor_current = float32(3493258) // 60Hz
	} else {
		rms_factor_current = float32(4191910) // 50Hz
	}


  dataAddress = make([]byte, 2)

	voltage_measure_1 := true
	voltage_measure_2 := true
	voltage_measure_3 := true

  for i:=0; i<=24; i++ {

    switch (i) {

      case 0:
        // current phase a
        dataAddress[0] = 0x43;//0x43C0 (AIRMS; Current rms an A)
        dataAddress[1] = 0xC0;
        data = make([]byte, 4)
      case 1:
        // current phase b
        dataAddress[0] = 0x43;//0x43C2 (BIRMS; Current rms an B)
        dataAddress[1] = 0xC2;
        data = make([]byte, 4)
      case 2:
        // current phase c
        dataAddress[0] = 0x43;//0x43C4 (CIRMS; Current rms an C)
        dataAddress[1] = 0xC4;
        data = make([]byte, 4)
      case 3:
        // current n
        dataAddress[0] = 0x43;//0x43C6 (NIRMS; Current rms neutral conductor)
        dataAddress[1] = 0xC6;
        data = make([]byte, 4)
      case 4:
        // voltage phase a
        dataAddress[0] = 0x43;//0x43C1 (AVRMS; Voltage rms an A)
        dataAddress[1] = 0xC1;
        data = make([]byte, 4)
      case 5:
        // voltage phase b
        dataAddress[0] = 0x43;//0x43C3 (BVRMS; Voltage rms an B)
        dataAddress[1] = 0xC3;
        data = make([]byte, 4)
      case 6:
        // voltage phase c
        dataAddress[0] = 0x43;//0x43C5 (CVRMS; Voltage rms an C)
        dataAddress[1] = 0xC5;
        data = make([]byte, 4)
			case 7:
        //  Phase A total active power.
        dataAddress[0] = 0xE5;//0xE513 (AWATT total active power an A)
        dataAddress[1] = 0x13;
        data = make([]byte, 4)
      case 8:
        //  Phase A total active power.
        dataAddress[0] = 0xE5;//0xE514 (BWATT total active power an B)
        dataAddress[1] = 0x14;
        data = make([]byte, 4)
      case 9:
        //  Phase A total active power.
        dataAddress[0] = 0xE5;//0xE515 (CWATT total active power an C)
        dataAddress[1] = 0x15;
        data = make([]byte, 4)
      case 10:
        // cosphi phase a
        dataAddress[0] = 0xE6;//0xE601 (ANGLE0 cosphi an A)
        dataAddress[1] = 0x01;
        data = make([]byte, 2)
      case 11:
        // cosphi phase b
        dataAddress[0] = 0xE6;//0xE602 (ANGLE1 cosphi an B)
        dataAddress[1] = 0x02;
        data = make([]byte, 2)
      case 12:
        // cosphi phase c
        dataAddress[0] = 0xE6;//0xE603 (ANGLE1 cosphi an B)
        dataAddress[1] = 0x03;
        data = make([]byte, 2)
      case 13:
        // frequency phase a
        register := []byte {0xE7, 0x00, 0x1C} //MMODE-Register measure frequency at VA
        err := d.Write(register)
        if err != nil {
            panic(err)
        }
					time.Sleep(50 * time.Millisecond)
        dataAddress[0] = 0xE6;//0xE607 (PERIOD)
        dataAddress[1] = 0x07;
        data = make([]byte, 2)
      case 14:
        // frequency phase b
        register := []byte {0xE7, 0x00, 0x1D} //MMODE-Register measure frequency at VB
        err = d.Write(register)
        if err != nil {
            panic(err)
        }
					time.Sleep(50 * time.Millisecond)
        dataAddress[0] = 0xE6;//0xE607 (PERIOD)
        dataAddress[1] = 0x07;
        data = make([]byte, 2)
      case 15:
        // frequency phase c
        register := []byte {0xE7, 0x00, 0x1E} //MMODE-Register measure frequency at VC
        err = d.Write(register)
        if err != nil {
            panic(err)
        }
				time.Sleep(50 * time.Millisecond)
        dataAddress[0] = 0xE6;//0xE607 (PERIOD)
        dataAddress[1] = 0x07;
        data = make([]byte, 2)
			case 16:
        //  Phase A total apparent power.
        dataAddress[0] = 0xE5;//0xE519 (AVA total apparent power an A)
        dataAddress[1] = 0x19;
        data = make([]byte, 4)
      case 17:
        //  Phase A total apparent power.
        dataAddress[0] = 0xE5;//0xE51A (BVA total apparent power an B)
        dataAddress[1] = 0x1A;
        data = make([]byte, 4)
      case 18:
        //  Phase A total apparent power.
        dataAddress[0] = 0xE5;//0xE51B (CVA total apparent power an C)
        dataAddress[1] = 0x1B;
        data = make([]byte, 4)
			case 19:
				//  Phase A total reactive power.
				dataAddress[0] = 0xE5;//0xE516 (AVAR total reactive power an A)
				dataAddress[1] = 0x16;
				data = make([]byte, 4)
			case 20:
				//  Phase A total reactive power.
				dataAddress[0] = 0xE5;//0xE517 (BVAR total reactive power an B)
				dataAddress[1] = 0x17;
				data = make([]byte, 4)
			case 21:
				//  Phase A total reactive power.
				dataAddress[0] = 0xE5;//0xE518 (CVAR total reactive power an C)
				dataAddress[1] = 0x18;
				data = make([]byte, 4)
    }

    // for j:=0; j<SAMPLES; j++ {

      err = d.Write(dataAddress)
      if err != nil {
          panic(err)
      }
      err = d.Read(data)
      if err != nil {
          panic(err)
      }

      switch (i) {
      case 0, 1, 2, 3, 4, 5, 6:
          // outcome = outcome + float32(FACTOR_3*int(data[0])+FACTOR_2*int(data[1])+FACTOR_1*int(data[2])+int(data[3]))
					outcome = float32(FACTOR_2*int(data[1])+FACTOR_1*int(data[2])+int(data[3]))
					fmt.Printf("I: %d, Data0: %x, Data1: %x, Data2: %x, Data3: %x \n",i, data[0], data[1], data[2], data[3])
      case 10 ,11 ,12, 13, 14, 15:
          // outcome = outcome + float32(FACTOR_1*int(data[0])+int(data[1]))
					outcome = float32(FACTOR_1*int(data[0])+int(data[1]))
		  case 7, 8, 9, 16, 17, 18, 19, 20, 21:
					// outcome = outcome + float32(FACTOR_3*int(data[0])+FACTOR_2*int(data[1])+FACTOR_1*int(data[2])+int(data[3]))
					outcome = float32(FACTOR_3*int(data[0])+FACTOR_2*int(data[1])+FACTOR_1*int(data[2])+int(data[3]))
      }


    // }

    // outcome = outcome / float32(SAMPLES)



    switch (i) {
      case 0:
        values[0] = ((((outcome * 0.3535) / rms_factor_current) / CURRENT_RESISTOR_A) / CURRENT_CLAMP_FACTOR_A) * 100.0 * OFFSET_CURRENT_A
      case 1:
        values[1] = ((((outcome * 0.3535) / rms_factor_current) / CURRENT_RESISTOR_B) / CURRENT_CLAMP_FACTOR_B) * 100.0 * OFFSET_CURRENT_B
      case 2:
        values[2] = ((((outcome * 0.3535) / rms_factor_current) / CURRENT_RESISTOR_C) / CURRENT_CLAMP_FACTOR_C) * 100.0 * OFFSET_CURRENT_C
      case 3:
        values[3] = ((((outcome * 0.3535) / rms_factor_current) / CURRENT_RESISTOR_N) / CURRENT_CLAMP_FACTOR_N) * 100.0 * OFFSET_CURRENT_N
      case 4:
				values[4] = float32(float32(outcome) / 1e+4)
				voltage_measure_1 = true
				if c.Measurevoltage1==0 || values[4] < 10 {
					values[4]= float32(c.Voltage1)
					voltage_measure_1 = false
				}
      case 5:
				values[5] = float32(float32(outcome) / 1e+4)
				voltage_measure_2 = true
				if c.Measurevoltage2==0 || values[5] < 10 {
					values[5]= float32(c.Voltage2)
					voltage_measure_2 = false
				}

      case 6:
				values[6] = float32(float32(outcome) / 1e+4)
				voltage_measure_3 = true
				if c.Measurevoltage3==0  || values[6] < 10 {
					values[6]= float32(c.Voltage3)
					voltage_measure_3 = false
				}

      case 7:
				values[i] = float32(outcome*POWER_CORRECTION_FACTOR_A)

				if c.Currentdirection1 == 1 {
					values[i] = values[i] * -1
				}

				if !voltage_measure_1 {
					values[7] = values[0] * values[4]
				}

			case 8:
				values[i] = float32(outcome*POWER_CORRECTION_FACTOR_B)

				if c.Currentdirection2 == 1 {
					values[i] = values[i] * -1
				}

				if !voltage_measure_2 {
					values[8] = values[1] * values[5]
				}

			case 9:
				values[i] = float32(outcome*POWER_CORRECTION_FACTOR_C)

				if c.Currentdirection3 == 1 {
					values[i] = values[i] * -1
				}

				if !voltage_measure_3 {
					values[9] = values[2] * values[6]
				}


      case 10:
        values[10] = float32(math.Cos(float64(outcome * FACTOR_CIRCLE * float32(c.Powerfrequency) / ADE7878_CLOCK * VAL)))

				if c.Currentdirection1 == 1 {
					values[i] = values[i] * -1
				}

				if c.Measurevoltage1==0 {
					values[10] = 1.0
				}

      case 11:
        values[11] = float32(math.Cos(float64(outcome * FACTOR_CIRCLE * float32(c.Powerfrequency) / ADE7878_CLOCK * VAL)))

				if c.Currentdirection2 == 1 {
					values[i] = values[i] * -1
				}

				if c.Measurevoltage2==0 {
					values[11] = 1.0
				}

      case 12:
        values[12] = float32(math.Cos(float64(outcome * FACTOR_CIRCLE * float32(c.Powerfrequency) / ADE7878_CLOCK * VAL)))

				if c.Currentdirection3 == 1 {
					values[i] = values[i] * -1
				}

				if c.Measurevoltage3==0 {
					values[12] = 1.0
				}

      case 13, 14, 15:
        values[i] = float32(ADE7878_CLOCK / (outcome+1))
			case 16 ,17 ,18:
				values[i] = float32(outcome)
			case 19:
				values[i] = float32(outcome)
				if c.Currentdirection1 == 1 {
					values[i] = values[i] * -1
				}
			case 20:
				values[i] = float32(outcome)
				if c.Currentdirection2 == 1 {
					values[i] = values[i] * -1
				}
			case 21:
				values[i] = float32(outcome)
				if c.Currentdirection3 == 1 {
					values[i] = values[i] * -1
				}
			case 22:
				// ??????????????????????????????????????
				if math.Signbit(float64(values[19])) {
					values[i] = (values[7]/POWER_CORRECTION_FACTOR_A/values[16])
				} else {
					values[i] = -1 * (values[7]/POWER_CORRECTION_FACTOR_A/values[16])
				}
			case 23:
				if math.Signbit(float64(values[20])) {
					values[i] = (values[8]/POWER_CORRECTION_FACTOR_B/values[17])
				} else {
					values[i] = -1 * (values[8]/POWER_CORRECTION_FACTOR_B/values[17])
				}
			case 24:
				if math.Signbit(float64(values[21])) {
					values[i] = (values[9]/POWER_CORRECTION_FACTOR_C/values[18])
				} else {
					values[i] = -1 * (values[9]/POWER_CORRECTION_FACTOR_C/values[18])
				}
    }

  }
	fmt.Printf("I1: %g  I2: %g  I3: %g  I4: %g  V1: %g  V2: %g  V3: %g  P1: %g  P2: %g  P3: %g  COS1: %g  COS2: %g  COS3: %g  F1: %g  F2: %g  F3: %g  AVA: %g  BVA: %g  CVA: %g  AVAR: %g  BVAR: %g  CVAR: %g  PFA: %g  PFB: %g  PFC: %g  \n",values[0],values[1],values[2],values[3],values[4],values[5],values[6],values[7],values[8],values[9],values[10],values[11],values[12],values[13],values[14],values[15],values[16],values[17],values[18],values[19],values[20],values[21],values[22],values[23],values[24]);

  return values

}

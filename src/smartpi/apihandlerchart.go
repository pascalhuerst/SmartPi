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
/*
File: apihandlerchart.go
Description: Handels API requests for charts
*/



package smartpi

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "time"
    "log"
    "strconv"
    "strings"
    "fmt"
    "math"
)

var Configfile string


func ServeChartValues(w http.ResponseWriter, r *http.Request) {

  type tChartValue struct {
    Time string `json:"time"`
    Value float32 `json:"value"`
  }

  type tChartSerie struct {
    Key string `json:"key"`
    Values []tChartValue `json:"values"`
    }

  // type tChartSeries []tChartSerie

	var timeSeries []tChartSerie

  vars := mux.Vars(r)
  from := vars["fromDate"]
  to := vars["toDate"]
  phaseId := vars["phaseId"]
  valueId := vars["valueId"]



  config := NewConfig()


  location := time.Now().Location()

  end, err := time.Parse(time.RFC3339,to)
  if err != nil {
    log.Fatal(err)
  }
  end = end.In(location)
	start, err := time.ParseInLocation(time.RFC3339,from,location)
  if err != nil {
    log.Fatal(err)
  }
  start = start.In(location)

  if end.Before(start) {
    start = start.AddDate(0,0,-1)
  }


  export := make([]string, 0)

	for i:=1; i<=3; i++ {
    if strings.Contains(phaseId,strconv.Itoa(i)) {
      export = append(export, valueId+"_"+strconv.Itoa(i))
    }
  }

  fmt.Println("ReadChartData "+config.Databasedir+" "+start.Format(time.RFC3339)+ " "+end.Format(time.RFC3339))

  data := ReadChartData(config.Databasedir, start, end)

	// fmt.Printf("%v", export)

	for _,valueelement := range export {
		row := 0
		val := 0.0
		var values []tChartValue
		for _,dataelement := range data {
			ti := dataelement.Date
			switch valueelement {
			case "current_1":
				val = dataelement.Current_1
			case "current_2":
				val = dataelement.Current_2
			case "current_3":
				val = dataelement.Current_3
      case "current_4":
				val = dataelement.Current_4
      case "voltage_1":
        val = dataelement.Voltage_1
      case "voltage_2":
        val = dataelement.Voltage_2
      case "voltage_3":
        val = dataelement.Voltage_3
      case "power_1":
        val = dataelement.Power_1
      case "power_2":
        val = dataelement.Power_2
      case "power_3":
        val = dataelement.Power_3
      case "cosphi_1":
        val = dataelement.Cosphi_1
      case "cosphi_2":
        val = dataelement.Cosphi_2
      case "cosphi_3":
        val = dataelement.Cosphi_3
      case "frequency_1":
        val = dataelement.Frequency_1
      case "frequency_2":
        val = dataelement.Frequency_2
      case "frequency_3":
        val = dataelement.Frequency_3
      case "energy_pos_1":
        val = dataelement.Energy_pos_1
      case "energy_pos_2":
        val = dataelement.Energy_pos_2
      case "energy_pos_3":
        val = dataelement.Energy_pos_3
      case "energy_neg_1":
        val = dataelement.Energy_neg_1
      case "energy_neg_2":
        val = dataelement.Energy_neg_2
      case "energy_neg_3":
        val = dataelement.Energy_neg_3


			}

			if math.IsNaN(val) {
        val = 0.0
      }
      values = append(values, tChartValue{Time: ti.Format(time.RFC3339), Value: float32( val )})
      row++
		}
		// fmt.Println(strconv.Itoa(index)+" "+valueelement)
		timeSeries = append(timeSeries, tChartSerie{Key: valueelement, Values: values})
	}



  // JSON output of request
  if err := json.NewEncoder(w).Encode(timeSeries); err != nil {
     panic(err)
  }
}


func ServeDayValues(w http.ResponseWriter, r *http.Request) {

  type tChartValue struct {
    Time string `json:"time"`
    Value float32 `json:"value"`
  }

  type tChartSerie struct {
    Key string `json:"key"`
    Values []tChartValue `json:"values"`
    }

  // type tChartSeries []tChartSerie

	var timeSeries []tChartSerie

  vars := mux.Vars(r)
  from := vars["fromDate"]
  to := vars["toDate"]
  phaseId := vars["phaseId"]
  valueId := vars["valueId"]



  config := NewConfig()


  // location := time.Now().Location()

  end, err := time.Parse(time.RFC3339,to)
  if err != nil {
    log.Fatal(err)
  }
  // end = end.In(location)
	start, err := time.Parse(time.RFC3339,from)
  if err != nil {
    log.Fatal(err)
  }
  // start = start.In(location)

  if end.Before(start) {
    start = start.AddDate(0,0,-1)
  }


  export := make([]string, 0)

	for i:=1; i<=3; i++ {
    if strings.Contains(phaseId,strconv.Itoa(i)) {
      export = append(export, valueId+"_"+strconv.Itoa(i))
    }
  }

  fmt.Println("ReadDayData "+config.Databasedir+" "+start.Format(time.RFC3339)+ " "+end.Format(time.RFC3339)+" |"+start.Location().String()+"|| "+start.Local().String())

  data := ReadDayData(config.Databasedir, start, end)

	// fmt.Printf("%v", export)

	for _,valueelement := range export {
		row := 0
		val := 0.0
		var values []tChartValue
		for _,dataelement := range data {
			ti := dataelement.Date
      fmt.Println(ti.Format(time.RFC3339))
			switch valueelement {
			case "current_1":
				val = dataelement.Current_1
			case "current_2":
				val = dataelement.Current_2
			case "current_3":
				val = dataelement.Current_3
      case "current_4":
				val = dataelement.Current_4
      case "voltage_1":
        val = dataelement.Voltage_1
      case "voltage_2":
        val = dataelement.Voltage_2
      case "voltage_3":
        val = dataelement.Voltage_3
      case "power_1":
        val = dataelement.Power_1
      case "power_2":
        val = dataelement.Power_2
      case "power_3":
        val = dataelement.Power_3
      case "cosphi_1":
        val = dataelement.Cosphi_1
      case "cosphi_2":
        val = dataelement.Cosphi_2
      case "cosphi_3":
        val = dataelement.Cosphi_3
      case "frequency_1":
        val = dataelement.Frequency_1
      case "frequency_2":
        val = dataelement.Frequency_2
      case "frequency_3":
        val = dataelement.Frequency_3
      case "energy_pos_1":
        val = dataelement.Energy_pos_1
      case "energy_pos_2":
        val = dataelement.Energy_pos_2
      case "energy_pos_3":
        val = dataelement.Energy_pos_3
      case "energy_neg_1":
        val = dataelement.Energy_neg_1
      case "energy_neg_2":
        val = dataelement.Energy_neg_2
      case "energy_neg_3":
        val = dataelement.Energy_neg_3
			}

			if math.IsNaN(val) {
        val = 0.0
      }
      values = append(values, tChartValue{Time: ti.Format(time.RFC3339), Value: float32( val )})
      row++
		}
		// fmt.Println(strconv.Itoa(index)+" "+valueelement)
		timeSeries = append(timeSeries, tChartSerie{Key: valueelement, Values: values})
	}



  // JSON output of request
  if err := json.NewEncoder(w).Encode(timeSeries); err != nil {
     panic(err)
  }
}

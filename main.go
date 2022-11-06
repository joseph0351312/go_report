package main

import (
	"fmt"

	"image/color"
	"math"
	"os"
	//	"math/rand"
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vinta/pangu"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	//	"log"
	//	"encoding/json"
	"gonum.org/v1/plot/vg/draw"
)

type Summarys struct {
	Data_R []*Summary
	Data_W []*Summary
}
type Summary struct {
	Name        string
	Succ        float64
	Fail        float64
	Send_Rate   float64
	Max_Latency float64
	Min_Latency float64
	Avg_Latency float64
	Throughput  float64
}

type Docker_performance struct {
	CPU_max     float64
	CPU_avg     float64
	Memory_max  float64
	Memory_avg  float64
	Traffic_In  float64
	Traffic_Out float64
	Disc_Write  float64
	Disc_Read   float64
}

func getPoints(yy, xx []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xx))
	for i, _ := range yy {
		pts[i].X = xx[i]
		pts[i].Y = yy[i]
	}
	return pts
}

func plot_Summary(data *Summarys, sends []float64, increment float64, title []string) {
	var init_latencys, read_latencys []float64
	var init_latencys_str, read_latencys_str []string
	var max_r_Latency, max_w_Latency, min_r_Latency, min_w_Latency float64
	init_latencys_str = append(init_latencys_str, "Write Latency ")
	read_latencys_str = append(read_latencys_str, "Query Latency ")
	max_r_Latency = 0
	max_w_Latency = 0
	min_r_Latency = data.Data_R[0].Avg_Latency
	min_w_Latency = data.Data_W[0].Avg_Latency
	//var maxs,mins,avgs[]float64
	p := plot.New()
	max_num := len(data.Data_W)
	for ii := 0; ii < max_num; ii++ {

		//	fmt.Println(ss)
		init_latencys = append(init_latencys, data.Data_W[ii].Avg_Latency)
		read_latencys = append(read_latencys, data.Data_R[ii].Max_Latency)

		init_latencys_str = append(init_latencys_str, fmt.Sprintf("%.2f", data.Data_W[ii].Avg_Latency))
		read_latencys_str = append(read_latencys_str, fmt.Sprintf("%.2f", data.Data_R[ii].Avg_Latency))

		if max_w_Latency < data.Data_W[ii].Avg_Latency {
			max_w_Latency = data.Data_W[ii].Avg_Latency
		}
		if max_r_Latency < data.Data_R[ii].Avg_Latency {
			max_r_Latency = data.Data_R[ii].Avg_Latency
		}
		if min_w_Latency < data.Data_W[ii].Avg_Latency {
			min_w_Latency = data.Data_W[ii].Avg_Latency
		}
		if min_r_Latency < data.Data_R[ii].Avg_Latency {
			max_r_Latency = data.Data_R[ii].Avg_Latency
		}
	}
	fmt.Printf("Latency\n   write \n\tmax:%.2f\n\tmin:%.2f\n", max_w_Latency, min_w_Latency)
	fmt.Printf("   read \n\tmax:%.2f\n\tmin:%.2f\n\n", max_r_Latency, min_r_Latency)
	p.X.Min = 0.0
	p.X.Max = sends[max_num-1]
	p.Y.Min = 0.0
	ff, _ := os.Create("output/latencys.csv")
	csvwriter := csv.NewWriter(ff)
	csvwriter.Write(title)
	fmt.Println(init_latencys)
	fmt.Println(read_latencys)
	csvwriter.Write(init_latencys_str)
	csvwriter.Write(read_latencys_str)
	csvwriter.Flush()
	p.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var rr = make([]plot.Tick, 0)
		i := 0
		for val := min; val <= max; val = val + increment {
			if i%100 == 0 {
				rr = append(rr, plot.Tick{val, fmt.Sprint(i)})
			} else {
				rr = append(rr, plot.Tick{val, ""})
			}
			i += 50
		}
		if (i-50)%100 != 0 {
			rr = append(rr, plot.Tick{max + 50, fmt.Sprint(i)})
		}
		return rr
	})
	p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var rr = make([]plot.Tick, 0)

		for val := min; val <= max; val = val + 1 {
			rr = append(rr, plot.Tick{val, fmt.Sprintf("%1.0f", val)})
		}

		return rr
	})
	p.Title.Text = "     "
	p.Title.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.X.Label.Text = "Send Rate (TPS)"
	p.Y.Label.Text = "Latency"
	p.X.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.Y.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.X.Tick.Label.Font.Size = 16
	p.Y.Tick.Label.Font.Size = 16
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter

	init_line, init_point, err := plotter.NewLinePoints(getPoints(init_latencys, sends))
	init_line.LineStyle.Width = vg.Points(2)
	init_line.LineStyle.Color = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	init_point.Shape = draw.BoxGlyph{}
	init_point.Color = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	if err != nil {
		panic(err)
	}
	read_line, read_point, err := plotter.NewLinePoints(getPoints(read_latencys, sends))
	read_line.LineStyle.Width = vg.Points(2)
	read_line.LineStyle.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	read_point.Shape = draw.BoxGlyph{}
	read_point.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	p.Add(init_line, init_point, read_line, read_point)

	p.Legend.Add("Write Latency ", init_line, init_point)
	p.Legend.Add("Query Latency ", read_line, read_point)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter
	/*
		"Max Latency", getPoints(maxs),
		"Min Latency", getPoints(mins),
		"Avg Latency", getPoints(avgs))*/
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "output/Latency.png"); err != nil {
		panic(err)
	}

}

func plot_Throughput(data *Summarys, sends []string, increment float64, title []string) {
	var init_throughput, read_throughput []float64
	var init_throughput_str, read_throughput_str []string
	var max_w_throughput, max_r_throughput, min_w_throughput, min_r_throughput float64
	max_r_throughput = 0.0
	max_w_throughput = 0.0
	min_r_throughput = data.Data_R[0].Throughput
	min_w_throughput = data.Data_W[0].Throughput
	init_throughput_str = append(init_throughput_str, "Write Throughput (TPS)")
	read_throughput_str = append(read_throughput_str, "Query Throughput (TPS)")

	//var maxs,mins,avgs[]float64
	p := plot.New()
	for ii := 0; ii < len(data.Data_W); ii++ {
		init_throughput = append(init_throughput, data.Data_W[ii].Throughput)
		read_throughput = append(read_throughput, data.Data_R[ii].Throughput)

		init_throughput_str = append(init_throughput_str, fmt.Sprintf("%.f", data.Data_W[ii].Throughput))
		read_throughput_str = append(read_throughput_str, fmt.Sprintf("%.f", data.Data_R[ii].Throughput))

		if max_r_throughput < data.Data_R[ii].Throughput {
			max_r_throughput = data.Data_R[ii].Throughput
		}
		if max_w_throughput < data.Data_W[ii].Throughput {
			max_w_throughput = data.Data_W[ii].Throughput
		}
		if min_r_throughput > data.Data_R[ii].Throughput {
			min_r_throughput = data.Data_R[ii].Throughput
		}
		if min_w_throughput > data.Data_W[ii].Throughput {
			min_w_throughput = data.Data_W[ii].Throughput
		}
	}
	fmt.Printf("Throughput\n   write \n\tmax:%f\n\tmin:%f\n", max_w_throughput, min_w_throughput)
	fmt.Printf("   read \n\tmax:%f\n\tmin:%f\n", max_r_throughput, min_r_throughput)
	ff, _ := os.Create("output/throughput.csv")
	csvwriter := csv.NewWriter(ff)
	csvwriter.Write(title)
	csvwriter.Write(init_throughput_str)
	csvwriter.Write(read_throughput_str)
	csvwriter.Flush()
	w := vg.Points(20)
	init_bar := plotter.Values(init_throughput)
	read_bar := plotter.Values(read_throughput)
	barsA, err := plotter.NewBarChart(init_bar, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	barsB, err := plotter.NewBarChart(read_bar, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)
	p.X.Max = math.Floor((p.X.Max/10)+1.5) * 100
	p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var rr = make([]plot.Tick, 0)

		i := 0
		for val := min; val <= max; val = val + increment {
			if i%100 == 0 {
				rr = append(rr, plot.Tick{val, fmt.Sprint(i)})
			} else {
				rr = append(rr, plot.Tick{val, ""})
			}
			i += 50
		}
		return rr
	})
	p.Title.Text = "     "
	p.Title.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.NominalX(sends...)
	p.X.Label.Text = "Send Rate (TPS)"
	p.Y.Label.Text = "Throughput"
	p.X.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.Y.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.X.Tick.Label.Font.Size = 16
	p.Y.Tick.Label.Font.Size = 16
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.X.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.Y.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.Y.Min = 0.0
	p.Add(barsA, barsB)

	p.Legend.Add("Write Throughput (TPS)", barsA)
	p.Legend.Add("Query Throughput (TPS)", barsB)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.X.Tick.Label.Font.Size = 20
	p.Y.Tick.Label.Font.Size = 16

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "output/Throughput.png"); err != nil {
		panic(err)
	}

}

/*
func plot_docker(davg, dmax []*Docker_performance) {
	var cpu_max, cpu_max1, cpu_avg, cpu_avg1, memory_max, memory_avg []float64
	p := plot.New()
	for _, ss := range davg {
		cpu_max = append(cpu_max, ss.CPU_max)
		cpu_avg = append(cpu_avg, ss.CPU_avg)
		memory_max = append(memory_max, ss.Memory_max/10)
		memory_avg = append(memory_avg, ss.Memory_avg/10)
	}
	for _, ss := range dmax {
		cpu_max1 = append(cpu_max1, ss.CPU_max)
		cpu_avg1 = append(cpu_avg1, ss.CPU_avg)
	}
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{{1.0, " "}, {2.0, "2"}, {3.0, " "}, {4.0, " "}, {5.0, "5"}, {6.0, " "}, {7.0, " "}, {8.0, " "}, {9.0, " "}, {10.0, "10"}})
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{{10.0, "10"}, {20.0, "20"}, {30.0, "30"}, {40.0, "40"}, {50.0, "50"}, {60.0, "60"}, {70.0, "70"}, {80.0, "80"}, {90.0, "90"}, {100.0, "100"}})

	p.Title.Text = "Hyperledger Fabric Performance"
	p.X.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.Y.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.X.Label.Text = "Peer"
	p.Y.Label.Text = "CPU "
	p.X.Max = 12.0
	p.X.Min = 0.0
	p.X.Tick.Label.Font.Size = 16
	p.Y.Tick.Label.Font.Size = 16
	cmax, _ := plotter.NewLine(getPoints(cpu_max))
	cmax.LineStyle.Width = vg.Points(4.5)
	cmax.LineStyle.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	cavg, _ := plotter.NewLine(getPoints(cpu_avg))
	cavg.LineStyle.Width = vg.Points(4)
	cavg.LineStyle.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	mmax, _ := plotter.NewLine(getPoints(memory_max))
	mmax.LineStyle.Width = vg.Points(2.5)
	mmax.LineStyle.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	mavg, _ := plotter.NewLine(getPoints(memory_avg))
	mavg.LineStyle.Width = vg.Points(2)
	mavg.LineStyle.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}

	cmax1, _ := plotter.NewLine(getPoints(cpu_max1))
	cmax1.LineStyle.Width = vg.Points(4.5)
	cmax1.LineStyle.Color = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	cavg1, _ := plotter.NewLine(getPoints(cpu_avg1))
	cavg1.LineStyle.Width = vg.Points(4)
	cavg1.LineStyle.Color = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	p.Add(cmax, cavg, cmax1, cavg1)
	p.Legend.Add("Max node CPU max", cmax1)
	p.Legend.Add("Max node CPU avg ", cavg1)
	p.Legend.Add("Avg node CPU max", cmax)
	p.Legend.Add("Avg node CPU avg ", cavg)
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter

	/*p.Legend.Add("Memory max[10/MB]", mmax)
	p.Legend.Add("Memory avg[10/MB]", mavg)
	p.Add(plotter.NewGrid())
	p.Y.Max = 90.0
	p.Y.Min = 0.0
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "docker_all.png"); err != nil {
		panic(err)
	}
}*/
func get_max(ss []string) *Docker_performance {
	var a float64
	var aa Docker_performance
	a = 0.0
	for _, line_str := range ss {
		sss := strings.Split(line_str, "</td>")
		if len(sss) != 1 {
			sss[0] = strings.Replace(sss[0], ".com", "", -1)
			sss[0] = strings.Replace(sss[0], "peer", "", -1)
			//	fmt.Println(sss)
			cpu_Maxs, _ := strconv.ParseFloat(sss[1], 64)
			cpu_Avgs, _ := strconv.ParseFloat(sss[2], 64)
			memory_Maxs, _ := strconv.ParseFloat(sss[3], 64)
			memory_Avgs, _ := strconv.ParseFloat(sss[4], 64)
			if a < cpu_Maxs {
				a = cpu_Maxs
				aa = Docker_performance{CPU_max: cpu_Maxs, CPU_avg: cpu_Avgs, Memory_max: memory_Maxs, Memory_avg: memory_Avgs}
			}
		}
	}

	return &aa
}
func senads_scope(first, last, increment int64) []float64 {
	var a []float64
	for i := float64(first); i <= float64(last); i += float64(increment) {
		a = append(a, i)
	}
	return a
}
func senads_scope_string(first, last, increment int64) []string {
	var a []string
	for i := first; i <= last; i += increment {
		a = append(a, fmt.Sprintln(i))
	}
	return a
}
func main() {
	err_num := 0
LABEL1:
	start()
LABEL2:
	var aa string
	fmt.Println("\n===============================\n")
	fmt.Println("輸入yes重複繪製/輸入no結束程式\n\n===============================")
	fmt.Scanln(&aa)
	if strings.EqualFold(aa, "no") || strings.EqualFold(aa, "n") {
		goto LABEL_End
	} else if strings.EqualFold(aa, "yes") || strings.EqualFold(aa, "y") {
		goto LABEL1
	} else if err_num > 1 {
		fmt.Println("輸入錯誤超過兩次強制結束\n")
		goto LABEL_End
	} else {
		fmt.Println("輸入錯誤\n")
		err_num++
		goto LABEL2
	}
LABEL_End:
}
func start() {
	var dataname, title []string
	var first, last, increment int64

	fmt.Println("請依序輸入起始值 測試間格 終點值(用空白鍵分開)")
	fmt.Scanf("%d%d%d\n", &first, &increment, &last)
	fmt.Printf("輸入的起始值為:%d\n輸入的測試間格:%d\n輸入的終點值%d\n", first, increment, last)
	var All_data [][]string

	title = append(title, "Send TPS")

	for i := first; i <= last; i = i + increment {
		file_name := fmt.Sprintf("tps_%d.html", i)
		dataname = append(dataname, fmt.Sprintf("%d_TPS", i))
		d, _ := get_data(file_name)
		All_data = append(All_data, d)

		title = append(title, fmt.Sprintf("%d", i))
	}

	/*       s4, docker_str := get_data("1_order_15_peer")
	         davg = append(davg, bar_char(docker_str, "1_order_15_peer"))
	         dmax = append(dmax, get_max(docker_str))
	*/
	c := csv_writer(All_data, dataname)
	plot_Throughput(c, senads_scope_string(first, last, increment), float64(increment), title)
	plot_Summary(c, senads_scope(first, last, increment), float64(increment), title)
	//	plot_docker(davg, dmax)
}

//func docker_data()
func get_data(performance string) ([]string, []string) {
	var ss []string
	var sss []string

	f, err := os.Open("./html/" + performance)

	if err != nil {
		fmt.Println(err)
	}
	doc, _ := goquery.NewDocumentFromReader(f)
	doc.Find("#benchmarksummary > table > tbody  ").Each(func(i int, selection *goquery.Selection) {
		//		ss := selection.Html()
		//	 	sss = strings.Split(ss,"")
		s, _ := selection.Html()
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "<tr><th>Name</th><th>Succ</th><th>Fail</th><th>SendRate(TPS)</th><th>MaxLatency(s)</th><th>MinLatency(s)</th><th>AvgLatency(s)</th><th>Throughput(TPS)</th></tr>", "", -1)
		s = strings.Replace(s, "<tr>", "", -1)
		s = strings.Replace(s, "<td>", "", -1)
		ss = strings.Split(s, "</tr>")

	})
	doc.Find("#" + performance /*+"> table >tbody"*/).Each(func(i int, se *goquery.Selection) {
		s, _ := se.Html()
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "<tr>", "", -1)
		s = strings.Replace(s, "<td>", "", -1)
		s = strings.Replace(s, "<table>", "", -1)
		s = strings.Replace(s, "<tbody>", "", -1)
		s = strings.Replace(s, "</tbody>", "", -1)
		s = strings.Replace(s, "<th>Name</th><th>CPU%(max)</th><th>CPU%(avg)</th><th>Memory(max)[MB]</th><th>Memory(avg)[MB]</th><th>TrafficIn[MB]</th><th>TrafficOut[MB]</th><th>DiscWrite[KB]</th><th>DiscRead[B]</th>", "", -1)
		s = strings.Replace(s, "<h4>Resourcemonitor:docker</h4>", "", -1)
		str := strings.Split(s, performance)
		strr := strings.Split(str[4], "</table>")
		sss = strings.Split(strr[1], "</tr>")
		fmt.Println(sss[1])
	})
	f.Close()
	return ss, sss
}

func csv_writer_(ss [][]string, name []string) *Summarys {
	var xxx [][]string
	var S Summarys
	ff, _ := os.Create("output/summary.csv")
	csvwriter := csv.NewWriter(ff)
	csvwriter.Write([]string{"name", "Succ", "Fail", "Send Rate (TPS)", "Max Latency (s)", "Min Latency (s)", "Avg Latency (s)", "Throughput (TPS)"})
	for i, line_strs := range ss {

		for _, line_str := range line_strs {
			sss := strings.Split(line_str, "</td>")
			if sss[0] != "" {
				sss[0] = name[i]
				csvwriter.Write(sss)
				xxx = append(xxx, sss)
				//lang, err := json.Marshal(sss)
				Succ, _ := strconv.ParseFloat(sss[1], 64)
				Fail, _ := strconv.ParseFloat(sss[2], 64)
				Send, _ := strconv.ParseFloat(sss[3], 64)

				Max, _ := strconv.ParseFloat(sss[4], 64)
				Min, _ := strconv.ParseFloat(sss[5], 64)
				Avg, _ := strconv.ParseFloat(sss[6], 64)
				Throughput, _ := strconv.ParseFloat(sss[7], 64)
				aa := Summary{Name: name[i],
					Succ:        Succ,
					Fail:        Fail,
					Send_Rate:   Send,
					Max_Latency: Max,
					Min_Latency: Min,
					Avg_Latency: Avg,
					Throughput:  Throughput}

				S.Data_R = append(S.Data_R, &aa)
			}

		}

	}
	csvwriter.Flush()

	//fmt.Println(xxx)
	return &S
}
func csv_writer(ss [][]string, name []string) *Summarys {
	var xxx [][]string
	var r [][]string
	var S Summarys

	ff, _ := os.Create("output/summary_W.csv")

	csvwriter := csv.NewWriter(ff)

	csvwriter.Write([]string{"name", "Succ", "Fail", "Send Rate (TPS)", "Max Latency (s)", "Min Latency (s)", "Avg Latency (s)", "Throughput (TPS)"})

	for i, line_strs := range ss {
		var (
			Succ, Fail, Send, Max, Min, Avg, Throughput float64
		)
		Send = 0.0
		Max = 0.0
		Min = 9999999.0
		Avg = 0.0
		Throughput = 0.0
		for ii, line_str := range line_strs {
			sss := strings.Split(line_str, "</td>")

			if sss[0] != "" {
				//

				//lang, err := json.Marshal(sss)
				Succ, _ = strconv.ParseFloat(sss[1], 64)
				Fail, _ = strconv.ParseFloat(sss[2], 64)
				Send, _ = strconv.ParseFloat(sss[3], 64)

				Max, _ = strconv.ParseFloat(sss[4], 64)

				Min, _ = strconv.ParseFloat(sss[5], 64)

				Avg, _ = strconv.ParseFloat(sss[6], 64)
				//Avg = (Avg + cum) / 2
				Throughput, _ = strconv.ParseFloat(sss[7], 64)
				aa := Summary{Name: name[i],
					Succ:        Succ,
					Fail:        Fail,
					Send_Rate:   Send,
					Max_Latency: Max,
					Min_Latency: Min,
					Avg_Latency: Avg,
					Throughput:  Throughput}
				if ii%2 == 0 {
					csvwriter.Write(sss)

					S.Data_W = append(S.Data_W, &aa)
				} else {
					r = append(r, sss)
					S.Data_R = append(S.Data_R, &aa)
				}
				xxx = append(xxx, sss)

			}

		}

	}
	csvwriter.Flush()
	ff, _ = os.Create("output/summary_R.csv")
	csvwriter = csv.NewWriter(ff)

	csvwriter.Write([]string{"name", "Succ", "Fail", "Send Rate (TPS)", "Max Latency (s)", "Min Latency (s)", "Avg Latency (s)", "Throughput (TPS)"})
	csvwriter.WriteAll(r)
	csvwriter.Flush()
	//fmt.Println(xxx)
	return &S
}
func one_node(ss []string, oo string) {
	var xxx [][]string

	for _, line_str := range ss {
		sss := strings.Split(line_str, "</td>")
		xxx = append(xxx, sss)

	}
}
func stringReplacer(text string) string {
	replacer := strings.NewReplacer("\n", "", " ", "")

	return pangu.SpacingText(replacer.Replace(text))
}
func bar_char(ss []string, oo string, size int) *Docker_performance {
	var cpu_max, cpu_avg, memory_max, memory_avg []float64
	var cpu_Maxs, cpu_Avgs, memory_Maxs, memory_Avgs, i float64
	i = 0

	cpu_Maxs = 0.0
	cpu_Avgs = 0.0
	memory_Maxs = 0.0
	memory_Avgs = 0.0
	var names []string
	//	fmt.Print(ss[1])
	for ii, line_str := range ss {
		fmt.Println(ii, "\n", line_str, "\n")

		sss := strings.Split(line_str, "</td>")
		if len(sss) != 1 {
			sss[0] = strings.Replace(sss[0], ".example.com", "", -1)
			sss[0] = strings.Replace(sss[0], "peer", "", -1)
			names = append(names, sss[0])
			fmt.Println("\n***********************\n", sss)
			cm, _ := strconv.ParseFloat(sss[1], 64)
			cpu_Maxs += cm
			cpu_max = append(cpu_max, cm)
			cm, _ = strconv.ParseFloat(sss[2], 64)
			cpu_Avgs += cm
			cpu_avg = append(cpu_avg, cm)
			cm, _ = strconv.ParseFloat(sss[3], 64)
			memory_Maxs += cm
			memory_max = append(memory_max, cm/10)
			cm, _ = strconv.ParseFloat(sss[4], 64)
			memory_Avgs += cm

			memory_avg = append(memory_avg, cm/10)
			i++
		}

	}
	//	fmt.Println(len(names))
	g_cpu_max := plotter.Values(cpu_max)
	g_cpu_avg := plotter.Values(cpu_avg)
	g_memory_max := plotter.Values(memory_max)
	g_memory_avg := plotter.Values(memory_avg)

	p := plot.New()

	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(20)

	barsA, err := plotter.NewBarChart(g_cpu_max, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	barsB, err := plotter.NewBarChart(g_cpu_avg, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)

	barsC, err := plotter.NewBarChart(g_memory_max, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)
	barsC.Offset = 2 * w
	barsD, err := plotter.NewBarChart(g_memory_avg, w)
	if err != nil {
		panic(err)
	}
	barsD.LineStyle.Width = vg.Length(0)
	barsD.Color = plotutil.Color(2)
	barsD.Offset = w
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{{10.0, "10"}, {20.0, "20"}, {30.0, "30"}, {40.0, "40"}, {50.0, "50"}, {60.0, "60"}, {70.0, "70"}, {80.0, "80"}, {90.0, "90"}, {100.0, "100"}})

	p.X.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter
	p.Y.Label.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.Y.Max = 100.0
	p.Y.Min = 0.0
	p.Add(barsA, barsB, barsC, barsD)

	p.Legend.Add("CPU%(max)", barsA)
	p.Legend.Add("CPU%(avg)", barsB)
	p.Legend.Add("Memory(max) [10/MB]", barsC)
	p.Legend.Add("Memory(avg) [10/MB]", barsD)
	p.Legend.TextStyle.Font.Size = 0.7 * vg.Centimeter

	p.Legend.Top = true
	p.NominalX(names...)
	p.X.Tick.Label.Font.Size = 20
	p.Y.Tick.Label.Font.Size = 16
	var x_size vg.Length
	switch i {
	case 2:
		x_size = 10 * vg.Inch
		break
	case 5:
		x_size = 15 * vg.Inch
		break
	case 10:
		x_size = 22 * vg.Inch
		break
	}
	if err := p.Save(x_size, 10*vg.Inch, oo+"_docker.png"); err != nil {
		panic(err)
	}
	var aa = Docker_performance{CPU_max: cpu_Maxs / i, CPU_avg: cpu_Avgs / i, Memory_max: memory_Maxs / i, Memory_avg: memory_Avgs / i}
	return &aa
}

此工具為golang編寫，已經過編譯

需要把Caliper的測試資料放進html
依照tps的數值將html命名成tps_數字.html
ex tps為50  tps_50.html

執行main.exe或plot.exe
根據敘述輸入數值如:

請依序輸入起始值 測試間格 終點值(用空白鍵分開)
50 50 550


結束程式後會產生
四個csv檔 
	summary_R.csv	有所有寫入的數據
	summary_W.csv	有所有讀出的數據
	latencys.csv	有在讀寫中 所有平均延遲時間的數據
	throughput.csv	有在讀寫中 所有交易數量
兩張圖
	
	Latency.png
	Throughput.png

plot.m.txt為matlab版本的繪圖工具


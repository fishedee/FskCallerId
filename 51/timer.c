#include "timer.h"
#include <reg52.h>
//定时器工具库
TimeHandlerType globalTimeHandler;
void delay(unsigned int z){
	unsigned x,y;
	for( x = z;x>0;x--){
		for( y = 110;y>0;y--){
			;
		}
	}
}
void TimerInit(){
	TMOD = ((TMOD&0xf0)|0x01);//定时器0的工作模式
	TH0=(65536-45872)/256; //设置定时器50ms触发一次，在11.0592d的晶振下
	TL0=(65536-45872)%256;
}
void TimerRun(TimeHandlerType timeHandler ){
	globalTimeHandler = timeHandler;
	EA=1;//开启全局中断
	TR0=1;//开启定时器
	ET0=1;//开启定时器中断
}
void T0_time()interrupt 1{
	TH0=(65536-45872)/256;
	TL0=(65536-45872)%256;
	globalTimeHandler();
}

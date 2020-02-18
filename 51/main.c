#include <reg52.h>
#include <stdio.h>
#include "timer.h"
#include "serial.h"
#include "lcd1602.h"


//业务函数
unsigned char num = 0;
unsigned char flag = 0;
char buffer[20];
void timeTick50ms(){
	num++;
	if(num==20){
		num=0;
		flag=1;
	}
}
void main(){
	unsigned int counter;
	int length;

	//模块初始化
	TimerInit();
	SerialInit();
	Lcd1602Init();

	//启动模块
   	TimerRun(timeTick50ms);
	SerialRun();

	//启动业务
	counter =0;
	while(1){
		if( flag == 1){
			flag = 0;
			counter++;
			length = sprintf(buffer,"counter: %d",counter);
			SerialSend(buffer,length);

			Lcd1602Write(1,buffer,length);
			length = sprintf(buffer,"counter2: %d",counter+1);
			Lcd1602Write(2,buffer,length);
		}
	}
}
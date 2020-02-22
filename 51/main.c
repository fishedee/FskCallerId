#include <reg52.h>
#include <stdio.h>
#include "timer.h"
#include "serial.h"
#include "lcd1602.h"


//业务函数
sbit key1 = P3^4;
sbit key2 = P3^5;
char flag = 0;
char buffer[64];

void checkKey(){
	flag = 0;
	if( key1 == 0 ){
		delay(10);
		if( key1 == 0){
			while(!key1);
			flag = 1;
		}
	}

	if( key2 == 0 ){
		delay(10);
		if( key2 == 0){
			while(!key2);
			flag = 2;
		}
	}
}
void main(){
	unsigned int counter;
	int length;

	//模块初始化
	SerialInit();
	//Lcd1602Init();

	//启动模块
	SerialRun();

	//启动业务
	while(1){
		checkKey();
		if( flag == 1){
			length = sprintf(buffer,"ATRING");
			SerialSend(buffer,length);
			//Lcd1602Write(1,buffer,length);
		}else if( flag == 2){
			length = sprintf(buffer,"ATF 0208091115018749403");
			SerialSend(buffer,length);
			//Lcd1602Write(1,buffer,length);
		}
	}
}
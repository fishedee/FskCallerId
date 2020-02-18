#include <reg52.h>
#include <stdio.h>
#include "timer.h"
#include "serial.h"
#include "lcd1602.h"


//ҵ����
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

	//ģ���ʼ��
	TimerInit();
	SerialInit();
	Lcd1602Init();

	//����ģ��
   	TimerRun(timeTick50ms);
	SerialRun();

	//����ҵ��
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
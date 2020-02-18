#include <reg52.h>
//串口工具库
void SerialInit(){
	TMOD =((TMOD&0x0f)|0x20);//定时器1的工作模式
	TH1=0xfd;//9600波特率
	TL1=0xfd;//9600波特率
	SM0=0; //设置串口为8位异步收发模式
	SM1=1;
	REN=0;//关闭串口接收功能
}

void SerialRun(){
	TR1=1;//启动T1定时器
	ES=0;//关闭串口的中断
}

void SerialSend(char* buffer, int length){
	unsigned int i;
	for( i = 0;i<length;i++){
		TI=0;
		SBUF=buffer[i];
		while(TI==0);
	}
	TI=0;
	SBUF='\n';
	while(TI==0);
}
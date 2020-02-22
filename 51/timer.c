#include "timer.h"
#include <reg52.h>
//��ʱ�����߿�
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
	TMOD = ((TMOD&0xf0)|0x01);//��ʱ��0�Ĺ���ģʽ
	TH0=(65536-45872)/256; //���ö�ʱ��50ms����һ�Σ���11.0592d�ľ�����
	TL0=(65536-45872)%256;
}
void TimerRun(TimeHandlerType timeHandler ){
	globalTimeHandler = timeHandler;
	EA=1;//����ȫ���ж�
	TR0=1;//������ʱ��
	ET0=1;//������ʱ���ж�
}
void T0_time()interrupt 1{
	TH0=(65536-45872)/256;
	TL0=(65536-45872)%256;
	globalTimeHandler();
}

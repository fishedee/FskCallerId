#include <reg52.h>
//LCD 1602���߿�
sbit lcden = P3^4;//Һ��ʹ�ܶ�
sbit lcdrs = P3^5;//Һ����������ѡ���
sfr lcddata = 0x80;//Һ�����ݴ��Ͷ�P0
void delay(unsigned int z){
	unsigned x,y;
	for( x = z;x>0;x--){
		for( y = 110;y>0;y--){
			;
		}
	}
}
void Lcd1602WriteCom(unsigned char com){
	lcdrs = 0;
	lcddata = com;
	delay(5);
	lcden = 1;
	delay(5);
	lcden = 0;
}

void Lcd1602WriteData(unsigned char inData){
	lcdrs = 1;
	lcddata = inData;
	delay(5);
	lcden = 1;
	delay(5);
	lcden = 0;
}
void Lcd1602Init(){
	Lcd1602WriteCom(0x38);//����Ϊ16x2��ʾ��5x7����8λ���ݽӿ�
	Lcd1602WriteCom(0x0c);//���ÿ���ʾ������ʾ���
	Lcd1602WriteCom(0x06);//дһ���ַ��󣬵�ַָ���1
	Lcd1602WriteCom(0x01);//��ʾ��0������ָ����0
}
void Lcd1602Write(unsigned char row, char* buffer, int length){
	unsigned char com;
	int i = 0;
	if( row == 1){
		com = 0x80;//��һ��д��
	}else{
		com = 0x80+0x40;//�ڶ���д��
	}
	Lcd1602WriteCom(com);
	for(i = 0;i< length ;i++){
		Lcd1602WriteData(buffer[i]);
		delay(5);
	}
}
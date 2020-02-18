#include <reg52.h>
//LCD 1602工具库
sbit lcden = P3^4;//液晶使能端
sbit lcdrs = P3^5;//液晶数据命令选择端
sfr lcddata = 0x80;//液晶数据传送端P0
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
	Lcd1602WriteCom(0x38);//设置为16x2显示，5x7点阵，8位数据接口
	Lcd1602WriteCom(0x0c);//设置开显示，不显示光标
	Lcd1602WriteCom(0x06);//写一个字符后，地址指针加1
	Lcd1602WriteCom(0x01);//显示清0，数据指针清0
}
void Lcd1602Write(unsigned char row, char* buffer, int length){
	unsigned char com;
	int i = 0;
	if( row == 1){
		com = 0x80;//第一行写入
	}else{
		com = 0x80+0x40;//第二行写入
	}
	Lcd1602WriteCom(com);
	for(i = 0;i< length ;i++){
		Lcd1602WriteData(buffer[i]);
		delay(5);
	}
}
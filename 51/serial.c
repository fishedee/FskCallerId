#include <reg52.h>
//���ڹ��߿�
void SerialInit(){
	TMOD =((TMOD&0x0f)|0x20);//��ʱ��1�Ĺ���ģʽ
	TH1=0xfd;//9600������
	TL1=0xfd;//9600������
	SM0=0; //���ô���Ϊ8λ�첽�շ�ģʽ
	SM1=1;
	REN=0;//�رմ��ڽ��չ���
}

void SerialRun(){
	TR1=1;//����T1��ʱ��
	ES=0;//�رմ��ڵ��ж�
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
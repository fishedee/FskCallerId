typedef void (*TimeHandlerType)();
void TimerInit();
void TimerRun(TimeHandlerType timeHandler );
void delay(unsigned int z);
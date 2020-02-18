typedef void (*TimeHandlerType)();
void TimerInit();
void TimerRun(TimeHandlerType timeHandler );
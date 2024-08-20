![go](https://img.shields.io/github/go-mod/go-version/pstrobl96/prusa_metrics_handler) 
![tag](https://img.shields.io/github/v/tag/pstrobl96/prusa_metrics_handler) 
![license](https://img.shields.io/github/license/pstrobl96/prusa_metrics_handler)

# Prusa Metrics Handler

This is simple utility, that takes metrics from printer, corrects timestamp and forward them into influx, easy as that.

Prusa 3D printers that are based on STM32 CPUs are unable to handle timestamp properly - they use delta timestamp - and you have to process them somewhere else.
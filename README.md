# Prusa Metrics Handler

This is simple utility, that takes metrics from printer, corrects timestamp and forward them into influx, easy as that.

Prusa 3D printers that are based on STM32 CPUs are unable to handle timestamp properly - they use delta timestamp - and you have to process them somewhere else.

If you want to send metrics into metrics processor then just run gcode bellow at your printer. Don't forget change the IP address to yours.

Details can be found in [Prusa_Firmware_Buddy](https://github.com/prusa3d/Prusa-Firmware-Buddy/blob/master/doc/metrics.md) repository.

```
M330 SYSLOG
M334 192.168.20.2 8514
M331 cpu_usage
M331 heap
M331 heap_free
M331 heap_total
M331 crash
M331 crash_stat
M331 crash_repeated
M331 excite_freq
M331 freq_gain
M331 tk_accel
M331 home_diff
M331 probe_z
M331 probe_z_diff
M331 probe_start
M331 probe_analysis
M331 heating_model_discrepancy
M331 esp_out
M331 esp_in
M331 fan_speed
M331 fan_hbr_speed
M331 ipos_x
M331 ipos_y
M331 ipos_z
M331 pos_x
M331 pos_y
M331 pos_z
M331 adj_z
M331 heater_enabled
M331 volt_bed_raw
M331 volt_bed
M331 volt_nozz_raw
M331 volt_nozz
M331 curr_nozz_raw
M331 curr_nozz
M331 curr_inp_raw
M331 curr_inp
M331 cur_mmu_imp
M331 oc_nozz
M331 oc_inp
M331 splitter_5V_current
M331 24VVoltage
M331 5VVoltage
M331 Sandwitch5VCurrent
M331 xlbuddy5VCurrent
M331 print_filename
M331 dwarf_board_temp
M331 dwarf_mcu_temp
M331 dwarfs_mcu_temp
M331 dwarfs_board_temp
M331 app_start
M331 maintask_loop
M331 fsensor_raw
M331 fsensor
M331 side_fsensor_raw
M331 side_fsensor
M331 nozzle_pwm
M331 bed_pwm
M331 loadcell
M331 loadcell_hp
M331 loadcell_xy
M331 loadcell_age
M331 loadcell_value
M331 power_panic
M331 crash_length
M331 usbh_err_cnt
M331 media_prefetched
M331 points_dropped
M331 probe_window
M331 eeprom_write
M331 tmc_sg_x
M331 tmc_sg_y
M331 tmc_sg_z
M331 tmc_sg_e
M331 tmc_write
M331 tmc_read
M331 fan
M331 print_fan_act
M331 hbr_fan_act
M331 gui_loop_dur
M331 g425_cen
M331 g425_off
M331 g425_rxy
M331 g425_xy
M331 g425_rz
M331 g425_z
M331 g425_xy_dev
M331 gcode
M331 loadcell_scale
M331 loadcell_threshold
M331 loadcell_threshold_cont
M331 loadcell_hysteresis
M331 mmu_comm
M331 dwarf_fast_refresh_delay
M331 dwarf_picked_raw
M331 dwarf_parked_raw
M331 dwarf_heat_curr
M331 bed_state
M331 bed_curr
M331 bedlet_state
M331 bedlet_temp
M331 bedlet_target
M331 bedlet_pwm
M331 bedlet_reg
M331 bedlet_curr
M331 bed_mcu_temp
M331 modbus_reqfail
M331 gui_loop_dur
```
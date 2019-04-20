![splash image](images/GoClock.png) 

# GoClock
another clock in Golang

goal is to check complexity against [Racket Clock](https://github.com/nodrygo/RktAlarmClock) and [Julia Clock](https://github.com/nodrygo/jAnalogAlarmClock) 

__Go__ is very basic language compared to Racket and Julia   
__Julia__ with [Luxor](https://github.com/JuliaGraphics/Luxor.jl) is a great help where Go need basic Cairo API
__Racket__ is easy to write (event when not using htdp)

# Compilation AOT
__Julia__ is very very slow to compile and the binaries are very big (>300Mo) including lot of shared lib     
__Julia__ start time is slow
__Julia__ portability between same OS but different processor seem problematic 
__Gtk__ not well finished    

__Racket__ is slow to compile and binary size is correct (~16Mo) with required shared lib        
__Racket__ start time is correct    
__Racket__  come with everything include 

__Go__ is fast to compile and binary size is correct (~16Mo) without shared lib      
__Go__ start time is fast  
__Gtk__ I have used [gotk3](https://github.com/gotk3/gotk3) not so bad but unfinished    

# my TOOLS
for __Julia__ I use Atom : _excellent_    
for __Racket__ I use DrRacket : _good_    
for __Go__ I use [LiteIDE X](https://liteide.org/en/) : _not bad_      

# TO DO    
alarm part   
right click menu  (set alarm / Resize / move / quit)  


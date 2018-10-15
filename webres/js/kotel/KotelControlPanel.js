Ext.define('KotelControlPanel', {
    extend: 'Ext.panel.Panel',
    initComponent: function() {
        this.title = 'Управление';
        this.border = true;
        this.frame = true;
        this.region = 'center';
        //this.height = 100;
        this.margins = '2 2 2 2';
        this.id = 'kotelControlPanel';
        this.collapsible = true;
        this.collapsed = false;
        this.resizable = false;
        this.bodyPadding = 10;
        this.bodyStyle = 'padding:10px;';

        this.currVal = ['25','60','11','1.8','e0','p0'];
        this.destVal = ['25','60','11','1.8','e0','p0'];
        this.stage = ['ot','vs','va','pr','es','es'];
        this.va = ['2','4','7','9','11','14'];
   
        this.mode = 0;

        this.curSatgeInx = 0;
        this.curVaInx = 4;
        
        this.command = '';
        
        this.initForm();

        KotelControlPanel.superclass.initComponent.apply(this, arguments);
    },
    initForm: function() {
        this.papa = this.initConfig().papa;
        
        this.task = { scope: this,
            run: function() { this.blink() },
            interval: 500
        };
        
        this.html = '<div id="cpBase">'
                +'<div class="led" id="led1" on="on"></div><div class="led" id="led2"></div>'
                +'<div id="dig1"></div><div id="dig2"></div><div id="dot"></div>'
                +'<div id="leftButt"></div><div id="rightButt"></div><div id="okButt"></div>'
                +'</div>';
       
        this.tbar = [
            Ext.create('Ext.Button', {text: 'Отменить', scope: this, disabled: false, id: 'execButt',
                style: 'background-position: bottom center;', 
                 handler: function() { this.command = ''; }
            }),
            '->', '-',
            Ext.create('Ext.Button', {text: 'Применить', scope: this, disabled: false, id: 'execButt',
                style: 'background-position: bottom center;', 
                handler: function() { this.sendCommands(); }
            })
        ];
        
        this.listeners = { scope: this,
            afterrender: function(){ 
                this.led1 = document.getElementById("led1");
                this.d1 = document.getElementById("dig1");
                this.d2 = document.getElementById("dig2");
                this.dot = document.getElementById("dot");
                this.bl = document.getElementById("leftButt");
                this.br = document.getElementById("rightButt");
                this.bo = document.getElementById("okButt");
                this.bl.onclick = this.leftClick;
                this.br.onclick = this.rightClick;
                this.bo.onclick = this.okClick;
                
                this.dispCurrentView();
            }
        };
    },
    leftClick: function(ev) {
        var cmp =Ext.getCmp('kotelControlPanel');
        if(cmp.mode==0) {
            if(cmp.curSatgeInx==0) {
                cmp.curSatgeInx = 3;
            } else
               cmp.curSatgeInx--; 
            
            cmp.setLedPosition(cmp.curSatgeInx);            
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='ot') {
            var val = cmp.destVal[cmp.curSatgeInx] == 'mm' ? 25 : parseInt(cmp.destVal[cmp.curSatgeInx]);
            val--;
            if(val<25) { val = 'mm'; }
            
            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='vs') {
            var val = cmp.destVal[cmp.curSatgeInx] == 'mm' ? 35 : parseInt(cmp.destVal[cmp.curSatgeInx]);
            val--;
            if(val<35) { val = 'mm'; }
            
            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='va') {
            cmp.curVaInx = cmp.curVaInx == 0 ? 0 : (--cmp.curVaInx);
            var val = cmp.va[cmp.curVaInx];            
            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        
        cmp.dispCurrentView();
            
        cmp.command += "L";
    },    
    rightClick: function(ev) {
        var cmp =Ext.getCmp('kotelControlPanel');
        if(cmp.mode==0) {
            if(cmp.curSatgeInx==3) {
                cmp.curSatgeInx = 0;
            } else
               cmp.curSatgeInx++;             
            cmp.setLedPosition(cmp.curSatgeInx);
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='ot') {
            var val = cmp.destVal[cmp.curSatgeInx] == 'mm' ? 24 : (cmp.destVal[cmp.curSatgeInx] == '85' ? 84 : parseInt(cmp.destVal[cmp.curSatgeInx])) ;
            val++;
            //if(val>99) { val = "mm"; }

            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='vs') {
            var val = cmp.destVal[cmp.curSatgeInx] == 'mm' ? 34 : (cmp.destVal[cmp.curSatgeInx] == '70' ? 69 : parseInt(cmp.destVal[cmp.curSatgeInx])) ;
            val++;
            //if(val>99) { val = "mm"; }

            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        if(cmp.mode==1 && cmp.stage[cmp.curSatgeInx]=='va') {
            cmp.curVaInx = cmp.curVaInx == 5 ? 5 : (++cmp.curVaInx);
            var val = cmp.va[cmp.curVaInx];            
            cmp.destVal[cmp.curSatgeInx] = val+'';
        }
        
        cmp.dispCurrentView();
        
        cmp.command += "R";
    },
    okClick: function(ev) {
        var cmp =Ext.getCmp('kotelControlPanel');
        if(cmp.stage[cmp.curSatgeInx]=='pr') return;
        if(cmp.mode==0) {
            cmp.mode=1;
            cmp.startBlink();
            cmp.display(cmp.destVal[cmp.curSatgeInx]);
        } else {
           cmp.mode=0; 
           cmp.stopBlink(); 
           cmp.setDest();
           cmp.display(cmp.currVal[cmp.curSatgeInx]);
        }
        
        //cmp.dispCurrentView();
        
        cmp.command += "M";
    },    
    dispCurrentView: function() {
        if(this.mode==0) this.display(this.currVal[this.curSatgeInx]);
        else this.display(this.destVal[this.curSatgeInx]);
    },
    display: function(val) {
        var d1 = '';
        var d2 = '';
        if(val.length>1) {
            d1 = val.substr(0,1);
            if(val.indexOf('.')>-1) {
                d2 = val.substr(2,1);
                this.dot.style.opacity=1;
            } else {
                d2 = val.substr(1,1);
                this.dot.style.opacity=0;
            }
            
        } else {
            d2 = val.substr(0,1);
        }
        
        this.d1.style.backgroundImage ='url("../img/kotel/dr'+d1+'.png")';
        this.d2.style.backgroundImage ='url("../img/kotel/dr'+d2+'.png")';
    },
    sendCommands: function() {
        var comm = this.command; 
        this.command = '';
        Ext.Ajax.request({
            url: '/api/setcomm', scope: this, method: 'GET',
            params: {comm: comm},
            success: function(response, opts) {
              var ansv = Ext.decode(response.responseText);
              if(ansv.success) { 
                Ext.Msg.show({
                    title: 'Отправка команды',
                    msg: 'Команда отправлена.',
                    buttons: Ext.Msg.OK,
                    icon: Ext.MessageBox.INFO,
                    fn: null
                });
                //this.command = '';
              } else error_mes('Ошибка', 'ErrorCode:'+ansv.error.errorCode+"; "+ansv.error.errorMessage);  
            },
            failure: function() { }
        });
    },
    setDest: function() {
        Ext.Ajax.request({
            url: '/api/setdestt', scope: this, method: 'GET',
            params: {desttp: this.destVal[0], destto: this.destVal[1], destkw: this.destVal[2]},
            success: function(response, opts) {
              var ansv = Ext.decode(response.responseText);
              if(ansv.success) {  

              } else error_mes('Ошибка', 'ErrorCode:'+ansv.error.errorCode+"; "+ansv.error.errorMessage);  
            },
            failure: function() { }
        });
    },
    blink: function() {
        if(this.led1.style.opacity==1) {
            this.led1.style.opacity = 0.2;
        } else {
            this.led1.style.opacity = 1;
        }
    },
    startBlink: function() {
        Ext.TaskManager.start(this.task);
    },
    stopBlink: function() {
        Ext.TaskManager.stop(this.task);
        this.led1.style.opacity=1;
    },
    setLedPosition: function(idx) {
        if(idx==5) idx=4;
        this.led1.style.left = (81+idx*141) + 'px';
    },
    setDestinationValuest: function(desttp, destto, destkw) {
        //if()
    }

});
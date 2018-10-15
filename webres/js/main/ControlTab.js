Ext.define('ControlTab', {
    extend: 'Ext.tab.Panel',   
    region: 'center',
    title: 'Панель управления',   
    //layout: 'border',
    border: false, frame: false,
    initComponent: function() {
      this.initPanel(); 
      this.items =  [ this.kotelTabPanel ];
      ControlTab.superclass.initComponent.apply(this, arguments);
    }, 
    initPanel: function() {
      this.kotelTabPanel   = Ext.create('KotelTabPanel',  {papa: this} );
      //this.devices = Ext.create('Devices',{papa: this});      
      //this.formPanel    = Ext.create('SearchFormPanel',{papa: this});
    }
});



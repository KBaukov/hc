Ext.define('MapsTab', {
    extend: 'Ext.tab.Panel',   
    region: 'center',
    title: 'Катры',   
    border: false, frame: false,
    mapsData: [],
    initComponent: function() {
      this.loadMaps(); 
      //this.items =  [ this.users, this.devices, this.maps ];
      //MapsTab.superclass.initComponent.apply(this, arguments);
    }, 
    loadMaps: function() {
        Ext.Ajax.request({
            url: '/api/maps', scope: this, method: 'POST',
            success: function(response, opts) {
              var ansv = Ext.decode(response.responseText);
              if(ansv.success) {  
                this.mapsData = ansv.data;
                this.createTabs();
              }    
            },
            failure: function() { this.unmask(); }
        });
    },
    createTabs: function() {
        var n = this.mapsData.length;
        var tabs = new Array(n);
        for(var i=0; i<n; i++) {
            tabs[i] = Ext.create('MapPanel',  {mapData: this.mapsData[i] } );
        }
        this.items = tabs;
        MapsTab.superclass.initComponent.apply(this, arguments);
    }
    
});
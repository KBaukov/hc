Ext.define('Devices', {
    extend: 'Ext.grid.Panel',
    initComponent: function() {
        this.title = 'Устройства';
        this.border = true;
        this.frame = false;
        this.region = 'center';
        this.height = 400;
        this.loadMask = true;
        this.emptyText  = 'Нет данных';
        this.margins = '2 2 2 2';
        this.id = 'devicesGrid';
        this.count = 0 ;
        this.collapsible = true;
        this.collapsed = false;
        this.resizable = true;
        
        this.initData();
        this.initColumns();

        Devices.superclass.initComponent.apply(this, arguments);
    },
    initColumns: function() {
      this.columns = [            
            {header: 'ID', align: 'left', width: 40, dataIndex: 'id', flex: 1 }, 
            {header: 'Тип устройства', align: 'left', width: 120, dataIndex: 'type' },
            {header: 'Имя', align: 'left', width: 160, dataIndex: 'name' },
            {header: 'IP адрес', align: 'left', width: 160, dataIndex: 'ip' },
            {header: 'Статус', align: 'left', width: 110, dataIndex: 'active_flag', renderer: this.statusRenderer },
            {header: 'Описание', align: 'left', width: 500, dataIndex: 'description' }
      ];
    },
    initData: function() {
      this.papa = this.initConfig().papa;
      this.store =  Ext.create('Ext.data.JsonStore', {
          storeId: 'devicesData', autoLoad: true,   
            proxy: {
                type: 'ajax',
                url: '/api/devices',
                reader: {
                    type: 'json',
                    root: 'data',
                    idProperty: 'id'
                }
            },
          fields: [
            {name: 'id'}, {name: 'type'}, {name: 'name'}, {name: 'ip'},
            {name: 'active_flag'}, {name: 'description'}
          ]//,
      });
    },
    loadData: function(subsId, start, cnt) { 
      this.mask();
      Ext.Ajax.request({
        url: '/api/users', scope: this, method: 'GET',
        params: {subscription_id: subsId, first: start, max: cnt },
        success: function(response, opts) {
          this.unmask();
          var ansv = Ext.decode(response.responseText);
          if(ansv.success) {  
            this.store.loadData(ansv.data);            
            this.count = this.store.count();
          } else error_mes('Ошибка', '!!!!!!!!!!!!!!!!!!!!!!!');  
        },
        failure: function() { this.unmask(); }
      });

    },
    statusRenderer: function(val) { 
        if(val=='Y') return 'активный'; 
        else return 'не активный'; 
    }
});


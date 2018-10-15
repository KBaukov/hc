Ext.define('Users', {
    extend: 'Ext.grid.Panel',
    initComponent: function() {
        this.title = 'Пользователи';
        this.border = true;
        this.frame = false;
        this.region = 'center';
        this.loadMask = true;
        this.emptyText  = 'Нет данных';
        this.margins = '2 2 2 2';
        this.id = 'UsersGrid';
        this.queryParams = {};
        
        this.initData();
        this.initColumns();
        

        Users.superclass.initComponent.apply(this, arguments);
    },
    initColumns: function() {
      this.columns = [
          {header: 'ID', align: 'left', width: 70, dataIndex: 'id' },
          {header: 'Тип пользователя', align: 'left', width: 110, dataIndex: 'user_type' },
          {header: 'Login', align: 'left', width: 100, dataIndex: 'login' },
          {header: 'Статус', align: 'left', width: 110, dataIndex: 'active_flag', renderer: this.statusRenderer },
          {header: 'Дата последнего визита', align: 'left', width: 200, dataIndex: 'last_visit', renderer: Ext.util.Format.dateRenderer('d.m.Y H:i:s') }
      ];
    },  
    initData: function() {
      
        this.papa = this.initConfig().papa;
        this.store =  Ext.create('Ext.data.JsonStore', {
            autoLoad: false,
            proxy: {
                type: 'ajax',
                url: '/api/users',
                reader: {
                    type: 'json',
                    root: 'data',
                    idProperty: 'id'
                }
            },
            fields: [
              {name: 'id', type: 'int'}, 
              {name: 'user_type'},
              {name: 'login'}, 
              //{name: 'pass'}, 
              {name: 'active_flag'},
              {name: 'last_visit'}
            ]
        });
        this.listeners = { scope: this,
              afterrender: function() {
                  this.loadData();
              }
        };
 
    },
    loadData: function() {     
      this.mask();
      Ext.Ajax.request({
        url: '/api/users', scope: this, method: 'GET',
        params: { },
        success: function(response, opts) {
          this.unmask();
          var ansv = Ext.decode(response.responseText);
          if(ansv.success) {  
            this.store.loadData(ansv.data);                       
          } else error_mes('Ошибка', 'ErrorCode:'+ansv.error.errorCode+"; "+ansv.error.errorMessage);  
        },
        failure: function() { this.unmask(); }
      });
    },
    statusRenderer: function(val) { 
        if(val=='Y') return 'активный'; 
        else return 'не активный'; 
    },
    dateRenderer: function(val) {
      return Ext.util.Format.date(val, 'd.m.Y H:i:s');
    }
});


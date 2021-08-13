function AppSettings(){
    function init(defaultValue){
        if (defaultValue){
            window.appSettings = defaultValue
        }
        else
        {
            window.appSettings = {
                api: "api/",
                pagelimit: 10
              }
        }
    }
    function setValue(key, value){
        if (!window.appSettings){
            throw 'Application settings are not initialized.'
        }        
        if (key === 'pagelimit'){
            window.appSettings[key] = parseInt(value) || 10;
        } else{
            window.appSettings[key] = value;
        }
    }
    function getValue(key){
        if (!window.appSettings){
            throw 'Application settings are not initialized.'
        }
        return window.appSettings[key];
    }
    function getAllValues(){
        return Object.assign({}, window.appSettings)
    }   
    return {
        init: init,
        setValue: setValue,
        getValue: getValue,
        getAllValues: getAllValues
    } 
}

export default AppSettings()
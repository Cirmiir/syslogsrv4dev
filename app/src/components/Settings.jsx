import { useState } from 'react';
import AppSettings from '../utils/AppSettings';

var element = document.getElementById("appsettings")
var settings = JSON.parse(element.dataset.settings)

export function Settings(){
    const [state, setState] = useState(AppSettings.getAllValues())
    function settingsChange(event){
        const obj = {}
        obj[event.target.name] = event.target.value;
        changeHandler(obj)

    }
    function changeHandler(e){
        setState( prevValues => {
          return { ...prevValues,...e}
        })
      }
    function saveSettings(){
        Object.keys(state).map((key) =>(
            AppSettings.setValue(key, state[key])
        ))
    }
    return(<div className="description"> 
    {Object.keys(settings).map((key) =>(<div>
        <div className="row">
        <div className="col-sm-2">{key}</div> : <div className="col-sm-2"><input className="form-control" readOnly disabled value={settings[key]}></input></div>
        </div>
    </div>))}

    {Object.keys(state).map((key) =>(<div>
        <div className="row">
        <div className="col-sm-2">{key}</div> :
            <div class="col-sm-2">
                <input className="form-control col-sm-2" name={key} onChange={settingsChange} value={state[key]}></input>
            </div>
        </div>
    </div>))}
    <button className="btn btn-success" onClick={saveSettings}>Save</button>
    </div>
    )
}

export default Settings
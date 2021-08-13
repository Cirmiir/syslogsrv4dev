import { BrowserRouter as Router, Route  } from "react-router-dom";

import React from 'react'
import Setting from './pages/Setting'
import About from './pages/About'
import Index from './pages/Index'
import AppSettings from './utils/AppSettings'

import 'react-minimal-side-navigation/lib/ReactMinimalSideNavigation.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import "gridjs/dist/theme/mermaid.css";
import 'react-date-range/dist/styles.css'; 
import 'react-date-range/dist/theme/default.css'; 

window.dataStorage = {
  _storage: new WeakMap(),
  put: function (element, key, obj) {
      if (!this._storage.has(element)) {
          this._storage.set(element, new Map());
      }
      this._storage.get(element).set(key, obj);
  },
  get: function (element, key) {
      return this._storage.get(element).get(key);
  },
  has: function (element, key) {
      return this._storage.has(element) && this._storage.get(element).has(key);
  },
  remove: function (element, key) {
      var ret = this._storage.get(element).delete(key);
      if (!this._storage.get(element).size === 0) {
          this._storage.delete(element);
      }
      return ret;
  }
}

AppSettings.init()

function App() {
  return (
    <div>
    <Router>
        <div>
          <Route exact path="/">
            <Index />
          </Route>
          <Route path="/about">
            <About />
          </Route>
          <Route path="/settings">
            <Setting />
          </Route>
        </div>
      </Router>
    </div>
  );
}


export default App;

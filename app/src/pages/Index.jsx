import AppSettings from '../utils/AppSettings'
import { FilterForm } from '../components/FilterForm'
import { Layout } from '../components/Layout'
import { Grid, _ } from 'gridjs-react';

import React, { useRef} from 'react'

export function Index() {
    const severityLevelDecode = {
      0:"Emergency",
      1:"Alert",
      2:"Critical",
      3:"Error",
      4:"Warning",
      5:"Notice",
      6:"Info",
      7:"Debug"
    }
  
    function cellhandler(e, cell){
      console.log(e)
    }
    const columns = ['Host', 'Date', 'AppName', 'Message', 'Facility', { 
      name: 'Severity',
      formatter: function(cell,row, table){
        var colorClass = "";
        if (cell <= 3){
          colorClass = "error"
        } else if (cell <= 5){
          colorClass = "warning"
        }
        return _(<div className={colorClass}>{severityLevelDecode[cell]}</div>)
      },
      onClick: cellhandler
    }]
    const grid = useRef(null);
    const filterRef = useRef(null);
  
    function reloadGrid(){
      if(grid.current){
  
        const currentGrid = grid.current.getInstance();
        currentGrid.forceRender();
    }
    }
  
  
    return (
      <Layout>
        <div>
          <FilterForm reload={reloadGrid} ref={filterRef}></FilterForm>
          <div className='grid-wrapper'>
          <Grid
            className='grid'
            ref={grid}
            columns={columns}
            pagination={{
              enabled: true,
              limit: window.appSettings.pagelimit,
              server: {
                body: (prev, page, limit) => {
                  prev = prev || {}
                  prev.Page = page
                  prev.Limit = limit
                  return JSON.stringify(prev)
                }
              }
            }}
            server={{
              url: AppSettings.getValue("api"),
              method: 'POST',
              body: {},
              data: (opts)=>{
                return new Promise((resolve, reject) => {
                  const xhttp = new XMLHttpRequest();
                  xhttp.onreadystatechange = function() {
                    if (this.readyState === 4) {
                      if (this.status === 200) {
                        const resp = JSON.parse(this.response);
                        resolve({
                          data: resp.Data.map(d => [d.Host, d.Date, d.AppName, d.Message, d.Facility, d.Severity]),
                          total: resp.TotalCount,
                        });
                      } else {
                        reject();
                      }
                    }
                  };
                  var parameters = {}
                  parameters = window.dataStorage.get(document.getElementById("filter-form"), "filters" , this).getFilters();     
                  const res = Object.assign({}, parameters, JSON.parse(opts.body));
                  xhttp.open(opts.method, opts.url, true);
                  xhttp.send(JSON.stringify(res));
                });
            }
            }}
          />
          </div>
        </div>
      </Layout>
    );
  }

  export default Index
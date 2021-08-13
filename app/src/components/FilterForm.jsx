import React, { Component, useRef } from 'react'
import Select from 'react-select'
import { DateRangePicker } from 'react-date-range';
import { format } from "date-fns";
import ReactDOM from 'react-dom';

var components = []

var element = document.getElementById("filters").childNodes;
var len = element.length;
if(len --) do {
    let item = element[len]
    let items = JSON.parse(item.dataset.items).map(function(val) {
      return {value: val.Key, label: val.Name }
  });
    components.push({
      Title: item.dataset.Title,
      Name: item.dataset.name,
      Items: items,
    })
} while(len --);

function formatDateRange(range){
  return `From: ${format(range.startDate, "MMMM do, yyyy")} To: ${format(range.endDate, "MMMM do, yyyy")} `
}

function extraOffset(date){
  let hoursDiff = date.getHours() - date.getTimezoneOffset() / 60;
  let minutesDiff = (date.getHours() - date.getTimezoneOffset()) % 60;
  date.setHours(hoursDiff);
  date.setMinutes(minutesDiff);

  return date;
}


function utcDate(date){
  return new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
}

const selectionRange = {
  startDate: utcDate(new Date()),
  endDate: utcDate(new Date()),
  key: 'selection',
}

export class FilterForm extends Component { 
  constructor(props) {
    super(props);
    this.reload = props.reload;
    this.filters = {
      StartDate: utcDate(selectionRange.startDate),
      EndDate:utcDate(selectionRange.endDate)
    };
    this.state = {
      showCalendar: false,
      showFiltres: false,
      dateRange: selectionRange
    }
    this.search = this.search.bind(this)
    this.handleSelect = this.handleSelect.bind(this)
    this.expandFilters = this.expandFilters.bind(this)
    this.handleClickOutside = this.handleClickOutside.bind(this)
    this.filterSelect = this.filterSelect.bind(this)
    window.filtersComponent = this;
    this.calendarRef = React.createRef()
    this.filterForm = React.createRef()
  }  
  search(e) {
    e.preventDefault();
    this.props.reload();
  }; 

  getFilters(){
    return this.filters;
  }

  handleClickOutside(event) {
    const domNode = ReactDOM.findDOMNode(this.calendarRef.current);

    if (this.state.showCalendar && (!domNode || !domNode.contains(event.target))) {
        this.changeHandler({showCalendar: false})
    }
  }

  componentDidMount() {
    document.addEventListener('click', this.handleClickOutside, true);
    window.dataStorage.put(this.filterForm.current, "filters" , this);
  }

  componentWillUnmount() {
      document.removeEventListener('click', this.handleClickOutside, true);
  }
  handleSelect(ranges){
    this.filters.StartDate = utcDate(extraOffset(ranges.selection.startDate));
    this.filters.EndDate = utcDate(extraOffset(ranges.selection.endDate));
    this.state.dateRange.startDate = ranges.selection.startDate;
    this.state.dateRange.endDate = ranges.selection.endDate;
  }
  filterSelect(e, par){
    this.filters[par.name] = e.map(function(item) {return item.value})
  }

  changeHandler(e){
    this.setState( prevValues => {
      return { ...prevValues,...e}
    })
  }

  expandFilters(e){
    this.changeHandler({showFiltres: !this.state.showFiltres})
  }

  render(){
    return <div id="filter-form" ref={this.filterForm} >
      <div className={`filter-section-title ${this.state.showFiltres ? "": "turned"}`} onClick={this.expandFilters}>Filters</div>
      <div className={`filter-section ${this.state.showFiltres ? "": "collapsed"}`} >
      <div className="row date-title">
        <div className="filter-control d-inline w-auto date-range-title">
          <div className="" onClick={() => this.changeHandler({showCalendar: !this.state.showCalendar})}>{formatDateRange(this.state.dateRange)}</div>
        </div>
      </div>
      <div className="row daterange-anchor">
        <DateRangePicker
        ref={this.calendarRef}
        className={`calendar ${this.state.showCalendar ? "" : "closed"}`} 
        ranges={[selectionRange]}
        onChange={this.handleSelect}
        showSelectionPreview={true}
        moveRangeOnFirstSelection={false}
        months={2}
        direction="horizontal"/>
        </div>
      <div className="row filter-row">
        {components.map((c, key) => (
        <Select 
          options={c.Items}
          isMulti
          name={c.Name}
          placeholder={c.Name}
          className="basic-multi-select filter-control from-control d-inline w-auto"
          classNamePrefix="select"
          onChange={this.filterSelect} />)
          )}

        <div className="search-btn d-inline w-auto">
          <button onClick={this.search} className="btn btn-outline-success d-inline w-auto">Search</button>
        </div>
      </div>
      </div>
    </div>
    
  };
}
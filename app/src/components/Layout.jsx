
import React from "react";
import { NavSideBar } from './NavSideBar'

export const Layout = ({children}) =>{
    return (<div>
            <NavSideBar/>
                <div className="content">
                        {children}
                </div>
        </div>);
}
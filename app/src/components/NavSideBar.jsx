import { useState } from 'react';
import { Navigation } from 'react-minimal-side-navigation';
import { useHistory, useLocation } from "react-router-dom";
import  React from 'react';
import { MdDehaze, MdSettings, MdHome, MdInfoOutline } from "react-icons/md";

export function NavSideBar(){
    const history = useHistory();
    const location = useLocation();
    const [IsLeftMenuOpened, setIsLeftMenuOpened] = useState(false);
    return (

    <React.Fragment>
      <div className='top-navbar'>
        <MdDehaze className="burger" onClick={() => setIsLeftMenuOpened(true)}/>
      </div>
      <div className={`nav-menu ${IsLeftMenuOpened ? "" : "closed"}`}>
        <div className={`left-menu` }>
          <Navigation activeItemId={location.pathname}
          onSelect={({ itemId }) => { history.push(itemId); }}
                  items={[
                    {
                      title: 'Logs',
                      itemId: '/',
                      elemBefore: () => <MdHome/>
                    },
                    {
                      title: 'About',
                      itemId: '/about',
                      elemBefore: () => <MdInfoOutline/>
                    },
                    {
                      title: 'Settings',
                      itemId: '/settings',
                      elemBefore: () => <MdSettings/>
                    },
                  ]}
                />
        </div>
        <div className={`background`} onClick={() => setIsLeftMenuOpened(false)}></div>  
      </div>        
    </React.Fragment>)
}
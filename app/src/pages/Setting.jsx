import { Layout } from '../components/Layout'
import { Settings } from '../components/Settings'

export function Setting() {
    return ( 
      <Layout>
        <div>
          <div className="row">
          <h2>Settings</h2>
          </div>
          <Settings />
        </div>
      </Layout>   
    );
  }

 export default Setting
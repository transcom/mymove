import React from 'react';

import magnifyingGlass from '../../../../node_modules/uswds/src/img/search.svg';
import LoginButton from 'shared/User/LoginButton';

import './index.css';

function QueueHeader() {
  return (
    <header role="banner" className="header">
      <div className="officeHeaderOne">
        <h3>admin.move.mil</h3>
      </div>
      <div className="officeHeaderTwo">
        <button>Queues</button>
      </div>
      <div className="officeHeaderThree">
        <div className="forms-shared">
          <form>
            <input type="text" placeholder="Search" name="search" />
          </form>
        </div>
        <form className="usa-form forms-shared">
          <select name="options" id="options">
            <option value>All</option>
            <option value="locator">Move Locator</option>
            <option value="edipi">DOD ID</option>
            <option value="customerName">Name (Last, First)</option>
            <option value="gbl">GBL</option>
            <option value="order">Orders</option>
            <option value="status">Status</option>
          </select>
        </form>
        <button>
          <img src={magnifyingGlass} alt="Search" />
        </button>
      </div>
      <div className="officeHeaderFour">
        <ul className="usa-nav-primary">
          <li>Welcome, Janelle</li>
          <li>
            <LoginButton />
          </li>
        </ul>
      </div>
    </header>
  );
}

export default QueueHeader;

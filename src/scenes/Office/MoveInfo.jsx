import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Route, Switch, Redirect } from 'react-router-dom';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faExclamationTriangle from '@fortawesome/fontawesome-free-solid/faExclamationTriangle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';

import './office.css';
import TextBoxWithEditLink from 'shared/TextBoxWithEditLink';

const BasicsTabContent = () => {
  return (
    <div>
      <div>
        <button>Approve Basics</button>
        <button>Troubleshoot</button>
        <button>Cancel Move</button>
      </div>
      <div>
        <h2>Customer Info</h2>
        <TextBoxWithEditLink />
        <h2>Backup Info</h2>
        <TextBoxWithEditLink />
        <div>
          <h2>Orders</h2>
          <div className="form-group">
            <form>
              <div className="within-form-group">
                <div className="form-column">
                  <label>Orders number</label>
                  <input type="text" />
                </div>
                <div className="form-column">
                  <label>Date issued</label>
                  <input type="text" />
                </div>
              </div>
              <div className="form-column">
                <label>Move type</label>
                <select name="move-type">
                  <option value="permanent-change-of-station">
                    Permanent Change of Station
                  </option>
                  <option value="separation">Separation</option>
                  <option value="retirement">Retirement</option>
                  <option value="local-move">Local Move</option>
                  <option value="tdy">Temporary Duty</option>
                  <option value="dependent-travel">Dependent Travel</option>
                  <option value="bluebark">Bluebark</option>
                  <option value="various">Various</option>
                </select>
              </div>
              <div className="form-column">
                <label>Orders type</label>
                <select name="orders-type">
                  <option value="shipment-of-hhg-permitted">
                    Shipment of HHG Permitted
                  </option>
                  <option value="pcs-with-tdy-en-route">
                    PCS with TDY En Route
                  </option>
                  <option value="shipment-of-hhg-restricted-or-prohibited">
                    Shipment of HHG Restricted or Prohibited
                  </option>
                  <option value="hhg-restricted-area-hhg-prohibited">
                    HHG Restricted Area - HHG Prohibited
                  </option>
                  <option value="course-of-instruction-20-weeks-or-more">
                    Course of Instruction 20 Weeks or More
                  </option>
                  <option value="shipment-of-hhg-prohibited-but-authorized-within-20-weeks">
                    Shipment of HHG Prohibited but Authorized within 20 Weeks
                  </option>
                  <option value="delayed-approval-20-weeks-or-more">
                    Delayed Approval 20 Weeks or More
                  </option>
                </select>
              </div>
              <div className="form-column">
                <label>Report by</label>
                <input type="date" name="report-by-date" />
              </div>
              <div className="form-column">
                <label>Current duty station</label>
                <input type="text" name="current-duty-station" />
              </div>
              <div className="form-column">
                <label>New duty station</label>
                <input type="text" name="new-duty-station" />
              </div>
              <div>
                <div className="form-column">
                  <b>Entitlements</b>
                  <label>Household goods</label>
                  <input type="number" name="household-goods-weight" /> lbs
                </div>
                <div className="form-column">
                  <label>Pro-gear</label>
                  <input type="number" name="pro-gear-weight" /> lbs
                </div>
                <div className="form-column">
                  <label>Spouse pro-gear</label>
                  <input type="number" name="spouse-pro-gear-weight" /> lbs
                </div>
                <div className="form-column">
                  <label>Short-term storage</label>
                  <input type="number" name="short-term-storage-days" /> days
                </div>
                <div className="form-column">
                  <label>Long-term storage</label>
                  <input type="number" name="long-term-storage-days" /> days
                </div>
                <div className="form-column">
                  <input type="checkbox" name="dependents-authorized" />Dependents
                  authorized
                </div>
              </div>
              <button>Cancel</button>
              <button>Save</button>
            </form>
          </div>
        </div>
        <div>
          <h2>Accounting</h2>
          <div className="form-group">
            <form>
              <div className="form-column">
                <label>
                  Dept. Indicator
                  <input type="text" />
                </label>
              </div>
              <div className="form-column">
                <label>
                  TAC
                  <input type="text" />
                </label>
              </div>
              <button>Cancel</button>
              <button>Save</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

const PPMTabContent = () => {
  return <div>PPM</div>;
};

export default class MoveInfo extends Component {
  render() {
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds Todo">
            <h1>Move Info: Johnson, Casey</h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole Todo">
            <ul className="move-info-header-meta">
              <li>ID# 3938593893</li>
              <li>
                (303) 936-8181
                <FontAwesomeIcon
                  className="icon"
                  icon={faPhone}
                  flip="horizontal"
                />
                <FontAwesomeIcon className="icon" icon={faComments} />
              </li>
              <li>Locator# ABC89</li>
              <li>KKFA to HAFC</li>
              <li>Requested Pickup 5/10/18</li>
            </ul>
          </div>
        </div>
        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-one-whole Todo">
            <p>Displaying move {this.props.match.params.moveID}.</p>

            <RoutedTabs startPathWith={this.props.match.url}>
              <NavTab to="/basics">
                <span className="title">Basics</span>
                <span className="status">
                  <FontAwesomeIcon className="icon" icon={faPlayCircle} />
                  Status Goes Here
                </span>
              </NavTab>
              <NavTab to="/ppm">
                <span className="title">PPM</span>
                <span className="status">
                  <FontAwesomeIcon
                    className="icon"
                    icon={faExclamationTriangle}
                  />
                  Status Goes Here
                </span>
              </NavTab>
            </RoutedTabs>

            <div className="tab-content">
              <Switch>
                <Route
                  exact
                  path={`${this.props.match.url}`}
                  render={() => (
                    <Redirect replace to={`${this.props.match.url}/basics`} />
                  )}
                />
                <Route
                  path={`${this.props.match.url}/basics`}
                  component={BasicsTabContent}
                />
                <Route
                  path={`${this.props.match.url}/ppm`}
                  component={PPMTabContent}
                />
              </Switch>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

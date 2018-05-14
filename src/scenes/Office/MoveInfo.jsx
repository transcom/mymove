import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Route, Switch, Redirect } from 'react-router-dom';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert'; // eslint-disable-line
import AccountingPanel from './AccountingPanel';
import { loadMoveDependencies } from './ducks.js';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';
import faExclamationTriangle from '@fortawesome/fontawesome-free-solid/faExclamationTriangle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';

import './office.css';

const BasicsTabContent = props => {
  return (
    <div>
      <div>
        <div>
          <h2>Customer Info</h2>
          <span className="fake-link">Edit</span>
          <br />
          <div className="form-column">
            <label>Title (optional)</label>
            <input type="text" name="title" />
          </div>
          <div className="form-column">
            <label>First name</label>
            <input type="text" name="first-name" />
          </div>
          <div className="form-column">
            <label>Middle name (optional)</label>
            <input type="text" name="middle-name" />
          </div>
          <div className="form-column">
            <label>Last name</label>
            <input type="text" name="last-name" />
          </div>
          <div className="form-column">
            <label>Suffix (optional)</label>
            <input type="text" name="name-suffix" />
          </div>
          <div className="form-column">
            <label>DoD ID</label>
            <input type="number" name="dod-id" />
          </div>
          <div className="form-column">
            <label>Branch</label>
            <select name="branch">
              <option value="army">Army</option>
              <option value="navy">Navy</option>
              <option value="air-force">Air Force</option>
              <option value="marines">Marines</option>
              <option value="coast-guard">Coast Guard</option>
            </select>
          </div>
          <div className="form-column">
            <label>Rank</label>
            <select name="rank">
              <option value="E-7">E-7</option>
              <option value="another-rank">Another rank</option>
              <option value="and-another-rank">And another rank</option>
            </select>
          </div>
          <div className="form-column">
            <b>Contact</b>
            <label>Phone</label>
            <input type="tel" name="contact-phone-number" />
          </div>
          <div className="form-column">
            <label>Alternate phone</label>
            <input type="tel" name="alternate-contact-phone-number" />
          </div>
          <div className="form-column">
            <label>Email</label>
            <input type="text" name="contact-email" />
          </div>
          <div className="form-column">
            <label>Preferred contact methods</label>
            <div>
              <input
                type="checkbox"
                id="phone-preference"
                name="preferred-contact-phone"
              />
              <label htmlFor="phone-preference">Phone</label>
            </div>
            <div>
              <input
                type="checkbox"
                id="text-preference"
                name="preferred-contact-text-message"
              />
              <label htmlFor="text-preference">Text message</label>
            </div>
            <div>
              <input
                type="checkbox"
                id="email-preference"
                name="preferred-contact-email"
              />
              <label htmlFor="email-preference">Email</label>
            </div>
          </div>
          <div className="form-column">
            <b>Current Residence Address</b>
            <label>Address 1</label>
            <input type="text" name="contact-address-1" />
          </div>
          <div className="form-column">
            <label>Address 2</label>
            <input type="text" name="contact-address-2" />
          </div>
          <div className="form-column">
            <label>City</label>
            <input type="text" name="contact-city" />
          </div>
          <div className="form-column">
            <label>State</label>
            <input type="text" name="contact-state" />
          </div>
          <div className="form-column">
            <label>Zip</label>
            <input type="number" name="contact-zip" />
          </div>
        </div>
        <div>
          <h2>Backup Info</h2>
          <span className="fake-link">Edit</span>
          <br />
          <form>
            <div className="form-column">
              <b>Backup Contact 1</b>
              <label>Name</label>
              <input type="text" name="backup-contact-1-name" />
            </div>
            <div className="form-column">
              <label>Phone</label>
              <input type="tel" name="backup-contact-1-phone" />
            </div>
            <div className="form-column">
              <label>Email (optional)</label>
              <input type="text" name="backup-contact-1-email" />
            </div>
            <div className="form-column">
              <b>Authorization</b>
              <input type="radio" name="authorization" value="none" />
              <label htmlFor="none">None</label>
              <input
                type="radio"
                name="authorization"
                value="letter-of-authorization"
              />
              <label htmlFor="letter-of-authorization">
                Letter of Authorization
              </label>
              <input
                type="radio"
                name="authorization"
                value="sign-for-pickup-delivery-only"
              />
              <label htmlFor="sign-for-pickup-delivery-only">
                Sign for pickup/delivery only
              </label>
            </div>
            <div className="form-column">
              <b>Backup Mailing Address</b>
              <label>Mailing address 1</label>
              <input type="text" name="backup-contact-1-mailing-address-1" />
            </div>
            <div className="form-column">
              <label>Mailing address 2</label>
              <input type="text" name="backup-contact-1-mailing-address-2" />
            </div>
            <div className="form-column">
              <label>City</label>
              <input type="text" name="backup-contact-1-city" />
            </div>
            <div className="form-column">
              <label>State</label>
              <input type="text" name="backup-contact-1-state" />
            </div>
          </form>
        </div>
        <div>
          <h2>Orders</h2>
          <span className="fake-link">Edit</span>
          <br />
          <div className="form-group">
            <form>
              <div className="within-form-group">
                <div className="form-column">
                  <label>Orders number</label>
                  <input type="text" name="orders-number" />
                </div>
                <div className="form-column">
                  <label>Date issued</label>
                  <input type="text" name="date-issued" />
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
                  <input
                    type="checkbox"
                    id="dependents-checkbox"
                    name="dependents-authorized"
                  />
                  <label htmlFor="dependents-checkbox">
                    Dependents authorized
                  </label>
                </div>
              </div>
              <button>Cancel</button>
              <button>Save</button>
            </form>
          </div>
        </div>
        <AccountingPanel moveId={props.match.params.moveId} />
      </div>
    </div>
  );
};

const PPMTabContent = () => {
  return <div>PPM</div>;
};

class MoveInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }

  render() {
    // TODO: If the following vars are not used to load data, remove them.
    // const officeMove = this.props.officeMove || {};
    // const officeOrders = this.props.officeOrders || {};
    const officeServiceMember = this.props.officeServiceMember || {};
    // const officeBackupContacts = this.props.officeBackupContacts || []
    const officePPMs = this.props.officePPMs || [];

    if (
      !this.props.loadDependenciesHasSuccess &&
      !this.props.loadDependenciesHasError
    )
      return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="usa-grid">
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              Something went wrong contacting the server.
            </Alert>
          </div>
        </div>
      );
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>
              Move Info: {officeServiceMember.last_name},{' '}
              {officeServiceMember.first_name}
            </h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta">
              <li>ID# {officeServiceMember.id}</li>
              <li>
                {officeServiceMember.telephone}
                {officeServiceMember.phone_is_preferred && (
                  <FontAwesomeIcon
                    className="icon"
                    icon={faPhone}
                    flip="horizontal"
                  />
                )}
                {officeServiceMember.text_message_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faComments} />
                )}
                {officeServiceMember.email_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faEmail} />
                )}
              </li>
              <li className="Todo">Locator# ABC89</li>
              <li>KKFA to HAFC</li>
              <li>
                Requested Pickup {get(officePPMs, '[0].planned_move_date')}
              </li>
            </ul>
          </div>
        </div>

        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-three-fourths">
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
                  path={`${this.props.match.path}/basics`}
                  component={BasicsTabContent}
                />
                <Route
                  path={`${this.props.match.path}/ppm`}
                  component={PPMTabContent}
                />
              </Switch>
            </div>
          </div>
          <div className="usa-width-one-fourths">
            <div>
              <button>Approve Basics</button>
              <button>Troubleshoot</button>
              <button>Cancel Move</button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

MoveInfo.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  officeMove: state.office.officeMove,
  officeOrders: state.office.officeOrders,
  officeServiceMember: state.office.officeServiceMember,
  officeBackupContacts: state.office.officeBackupContacts,
  officePPMs: state.office.officePPMs,
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(MoveInfo);

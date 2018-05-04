import React, { Component } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { loadPpm } from 'scenes/Moves/Ppm/ducks';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import ppmBlack from 'shared/icon/ppm-black.svg';
import { indexBackupContacts } from 'scenes/ServiceMembers/ducks';

import './Review.css';

export class Review extends Component {
  componentDidMount() {
    this.props.loadPpm(this.props.match.params.moveId);
  }
  componentDidUpdate() {
    const service_member = get(this.props.loggedInUser, 'service_member');
    if (service_member) {
      this.props.indexBackupContacts(service_member.id);
    }
  }
  render() {
    const {
      pages,
      pageKey,
      currentPpm,
      currentBackupContacts,
      loggedInUser,
    } = this.props;
    const backupContact = currentBackupContacts;
    function getFullName() {
      const service_member = get(loggedInUser, 'service_member');
      if (!service_member) return;
      return `${service_member.first_name} ${service_member.middle_name ||
        ''} ${service_member.last_name} ${service_member.suffix || ''}`;
    }
    function getFullResAddress() {
      const residential_address = get(
        loggedInUser,
        'service_member.residential_address',
      );
      if (residential_address) {
        return `${
          residential_address.street_address_1
        } ${residential_address.street_address_2 || ''} ${
          residential_address.city
        } ${residential_address.state} ${residential_address.postal_code}`;
      }
    }
    function getFullBackupAddress() {
      const backup_mailing_address = get(
        loggedInUser,
        'service_member.backup_mailing_address',
      );
      if (backup_mailing_address) {
        return `${
          backup_mailing_address.street_address_1
        } ${backup_mailing_address.street_address_2 || ''} ${
          backup_mailing_address.city
        } ${backup_mailing_address.state} ${
          backup_mailing_address.postal_code
        }`;
      }
    }
    return (
      <WizardPage
        handleSubmit={no_op}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={true}
      >
        <h1>Review</h1>
        <p>
          You're almost done! Please review your details before we finalize the
          move.
        </p>
        <h3>Profile and Orders</h3>

        <div className="usa-grid-full review-content">
          <div className="usa-width-one-half review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Profile{' '}
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>{getFullName()}</td>
                </tr>
                <tr>
                  <td>Branch:</td>
                  <td>{get(loggedInUser, 'service_member.affiliation')}</td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td>{get(loggedInUser, 'service_member.rank')}</td>
                </tr>
                <tr>
                  <td> DoD ID#: </td>
                  <td>{get(loggedInUser, 'service_member.edipi')}</td>
                </tr>
                <tr>
                  <td> Current Duty Station: </td>
                  <td>
                    {get(loggedInUser, 'service_member.current_station.name')}
                  </td>
                </tr>
              </tbody>
            </table>

            <table className="review-Todo Todo">
              <tbody>
                <tr>
                  <th className="Todo">
                    Orders{' '}
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Orders Type: </td>
                  <td>Permanent Change of Station</td>
                </tr>
                <tr>
                  <td> Orders Date: </td>
                  <td> 06/01/2018 </td>
                </tr>
                <tr>
                  <td> Report-by Date:: </td>
                  <td> 07/11/2018</td>
                </tr>
                <tr>
                  <td> New Duty Station: </td>
                  <td> Fort Carson </td>
                </tr>
                <tr>
                  <td> Dependents? </td>
                  <td> Yes </td>
                </tr>
                <tr>
                  <td> Orders Uploaded: </td>
                  <td> 8 photos uploaded</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div className="usa-width-one-half review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Contact Info{' '}
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Best Contact Phone: </td>
                  <td>{get(loggedInUser, 'service_member.telephone')}</td>
                </tr>
                <tr>
                  <td> Alt. Phone: </td>
                  <td>
                    {get(loggedInUser, 'service_member.secondary_telephone')}
                  </td>
                </tr>
                <tr>
                  <td> Personal Email: </td>
                  <td>{get(loggedInUser, 'service_member.personal_email')}</td>
                </tr>
                <tr>
                  <td> Preferred Contact Method: </td>
                  <td>
                    {get(loggedInUser, 'service_member.current_station.name')}
                  </td>
                </tr>
                <tr>
                  <td> Current Mailing Address: </td>
                  <td>{getFullResAddress()}</td>
                </tr>
                <tr>
                  <td> Backup Mailing Address: </td>
                  <td>{getFullBackupAddress()}</td>
                </tr>
              </tbody>
            </table>
            {currentBackupContacts.map(contact => (
              <table>
                <tbody>
                  <tr>
                    <th>
                      Backup Contact Info{' '}
                      <span className="align-right">
                        <a href="about:blank">Edit</a>
                      </span>
                    </th>
                  </tr>
                  <tr>
                    <td> Backup Contact: </td>
                    <td>
                      <p>{contact.name}</p>
                      <p>{contact.permission}</p>
                    </td>
                  </tr>
                  <tr>
                    <td> Email: </td>
                    <td> {contact.email} </td>
                  </tr>
                  <tr>
                    <td> Phone: </td>
                    <td> {contact.telephone}</td>
                  </tr>
                </tbody>
              </table>
            ))}
          </div>
        </div>

        <div className="usa-grid-full ppm-container">
          <h3>
            <img src={ppmBlack} alt="PPM shipment" /> Shipment - You move your
            stuff (PPM)
          </h3>
          <div className="usa-width-one-half review-section ppm-review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Dates & Locations
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td className="Todo"> Move Date: </td>
                  <td className="Todo">
                    {' '}
                    {currentPpm && currentPpm.planned_move_date}
                  </td>
                </tr>
                <tr>
                  <td> Pickup ZIP Code: </td>
                  <td> {currentPpm && currentPpm.pickup_zip}</td>
                </tr>
                <tr>
                  <td> Delivery ZIP Code: </td>
                  <td> {currentPpm && currentPpm.destination_zip}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div className="usa-width-one-half review-section ppm-review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Weight
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Estimated Weight: </td>
                  <td> {currentPpm && currentPpm.weight_estimate} lbs</td>
                </tr>
                <tr>
                  <td> Estimated PPM Incentive: </td>
                  <td className="Todo">
                    {' '}
                    {currentPpm && currentPpm.estimated_incentive}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </WizardPage>
    );
  }
}

Review.propTypes = {
  currentPpm: PropTypes.object,
  currentBackupContacts: PropTypes.array,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPpm, indexBackupContacts }, dispatch);
}

function mapStateToProps(state) {
  const props = {
    ...state.ppm,
    ...state.loggedInUser,
    currentBackupContacts: state.serviceMember.currentBackupContacts,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Review);

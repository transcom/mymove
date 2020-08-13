import React from 'react';
import { Link, withRouter } from 'react-router-dom';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import { formatDateSM } from 'shared/formatters';
import { getFullSMName } from 'utils/moveSetupFlow';
import Address from './Address';

import './Review.css';

function ServiceMemberSummary(props) {
  const {
    orders,
    serviceMember,
    schemaRank,
    schemaAffiliation,
    schemaOrdersType,
    moveIsApproved,
    editOrdersPath,
    uploads,
  } = props;

  const rootPath = `/moves/review`;
  const editProfilePath = rootPath + '/edit-profile';

  const yesNoMap = { true: 'Yes', false: 'No' };

  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h2>
            Profile
            <span className="edit-section-link">
              <Link to={editProfilePath} className="usa-link">
                Edit
              </Link>
            </span>
          </h2>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Name</th>
            <td>{getFullSMName(serviceMember)}</td>
          </tr>
          <tr>
            <th scope="row">Branch</th>
            <td>{get(schemaAffiliation['x-display-value'], get(serviceMember, 'affiliation'))}</td>
          </tr>
          <tr>
            <th scope="row">Rank</th>
            <td>{get(schemaRank['x-display-value'], get(serviceMember, 'rank'))}</td>
          </tr>
          <tr>
            <th scope="row">DoD ID#</th>
            <td>{get(serviceMember, 'edipi')}</td>
          </tr>
          <tr>
            <th scope="row">Current duty station</th>
            <td>{get(serviceMember, 'current_station.name')}</td>
          </tr>
        </tbody>
      </table>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Contact info</th>
          </tr>
          <tr>
            <th scope="row">Best contact phone</th>
            <td>{get(serviceMember, 'telephone')}</td>
          </tr>
          <tr>
            <th scope="row">Personal email</th>
            <td>{get(serviceMember, 'personal_email')}</td>
          </tr>
          <tr>
            <th scope="row">Current mailing address</th>
            <td>
              <Address address={get(serviceMember, 'residential_address')} />
            </td>
          </tr>
        </tbody>
      </table>
      <div className="stackedtable-header">
        <div>
          <h2>
            Orders
            {moveIsApproved && '*'}
            {!moveIsApproved && (
              <span className="edit-section-link">
                <Link to={editOrdersPath} className="usa-link">
                  Edit
                </Link>
              </span>
            )}
          </h2>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Orders type</th>
            <td>{get(schemaOrdersType['x-display-value'], get(orders, 'orders_type'))}</td>
          </tr>
          <tr>
            <th scope="row">Orders date</th>
            <td>{formatDateSM(get(orders, 'issue_date'))}</td>
          </tr>
          <tr>
            <th scope="row">Report by date</th>
            <td>{formatDateSM(get(orders, 'report_by_date'))}</td>
          </tr>
          <tr>
            <th scope="row">New duty station</th>
            <td>{get(orders, 'new_duty_station.name')}</td>
          </tr>
          <tr>
            <th scope="row">Dependents</th>
            <td>{orders && yesNoMap[get(orders, 'has_dependents', '').toString()]}</td>
          </tr>
          <tr>
            <th scope="row">Orders</th>
            <td>{uploads && uploads.length}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

/*

    <div className="service-member-summary">
      <GridContainer>
        <Grid row>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <h2 className="heading">
                Profile
                <span className="edit-section-link">
                  <Link to={editProfilePath} className="usa-link">
                    Edit
                  </Link>
                </span>
              </h2>
              <table>
                <tbody>
                  <tr>
                    <td>Name </td>
                    <td>{getFullSMName(serviceMember)}</td>
                  </tr>
                  <tr>
                    <td>Branch </td>
                    <td>{get(schemaAffiliation['x-display-value'], get(serviceMember, 'affiliation'))}</td>
                  </tr>
                  <tr>
                    <td>Rank </td>
                    <td>{get(schemaRank['x-display-value'], get(serviceMember, 'rank'))}</td>
                  </tr>
                  <tr>
                    <td>DoD ID# </td>
                    <td>{get(serviceMember, 'edipi')}</td>
                  </tr>
                  <tr>
                    <td>Current duty station </td>
                    <td>{get(serviceMember, 'current_station.name')}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </Grid>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p>Contact info</p>
              <table>
                <tbody>
                  <tr>
                    <td>Best contact phone </td>
                    <td>{get(serviceMember, 'telephone')}</td>
                  </tr>
                  <tr>
                    <td>Personal email </td>
                    <td>{get(serviceMember, 'personal_email')}</td>
                  </tr>
                  <tr>
                    <td>Current mailing address </td>
                    <td>
                      <Address address={get(serviceMember, 'residential_address')} />
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </Grid>
        </Grid>
        <Grid row>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              {orders && (
                <Fragment>
                  <h2 className="heading">
                    Orders
                    {moveIsApproved && '*'}
                    {!moveIsApproved && (
                      <span className="edit-section-link">
                        <Link to={editOrdersPath} className="usa-link">
                          Edit
                        </Link>
                      </span>
                    )}
                  </h2>
                  <table>
                    <tbody>
                      <tr>
                        <td>Orders type </td>
                        <td>{get(schemaOrdersType['x-display-value'], get(orders, 'orders_type'))}</td>
                      </tr>
                      <tr>
                        <td>Orders date </td>
                        <td> {formatDateSM(get(orders, 'issue_date'))}</td>
                      </tr>
                      <tr>
                        <td>Report by date </td>
                        <td>{formatDateSM(get(orders, 'report_by_date'))}</td>
                      </tr>
                      <tr>
                        <td>New duty station </td>
                        <td> {get(orders, 'new_duty_station.name')}</td>
                      </tr>
                      <tr>
                        <td>Dependents </td>
                        <td> {orders && yesNoMap[get(orders, 'has_dependents', '').toString()]}</td>
                      </tr>
                      <tr>
                        <td>Orders </td>
                        <td>{uploads && uploads.length}</td>
                      </tr>
                    </tbody>
                  </table>
                </Fragment>
              )}
            </div>
          </Grid>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              {backupContacts.map((contact) => (
                <Fragment key={contact.id}>
                  <p className="heading">
                    Backup Contact Info
                    <span className="edit-section-link">
                      <Link to={editBackupContactPath} className="usa-link">
                        Edit
                      </Link>
                    </span>
                  </p>
                  <table key={contact.id}>
                    <tbody>
                      <tr>
                        <td> Backup Contact: </td>
                        <td>
                          {contact.name} <br />
                          { / * getFullBackupPermission(contact) * / }
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
                </Fragment>
              ))}
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
*/

ServiceMemberSummary.propTypes = {
  backupContacts: PropTypes.array.isRequired,
  serviceMember: PropTypes.object,
  schemaRank: PropTypes.object.isRequired,
  schemaAffiliation: PropTypes.object.isRequired,
  schemaOrdersType: PropTypes.object.isRequired,
  orders: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  editOrdersPath: PropTypes.string,
};

export default withRouter(ServiceMemberSummary);

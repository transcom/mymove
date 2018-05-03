import React, { Fragment, Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { loadPpm } from 'scenes/Moves/Ppm/ducks';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import './Review.css';

export class Review extends Component {
  componentDidMount() {
    this.props.loadPpm(this.props.match.params.moveId);
  }
  render() {
    const {
      pages,
      pageKey,
      currentServiceMember,
      currentPpm,
      currentBackupContacts,
    } = this.props;
    console.log(this.props);
    return (
      <WizardPage
        handleSubmit={no_op}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={true}
      >
        <Fragment>
          <h1>Review</h1>
          <p>
            You're almost done! Please review your details before we finalize
            the move.
          </p>
          <h3>Profile and Orders</h3>

          <div className="usa-grid-full review-content">
            <div className="usa-width-one-half review-section">
              <table>
                <tr>
                  <th>Profile Edit</th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>
                    Michael Owen Jones{' '}
                    {currentServiceMember && currentServiceMember.first_name}
                  </td>
                </tr>
                <tr>
                  <td> Branch: </td>
                  <td> Air Force </td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td> E-5 </td>
                </tr>
                <tr>
                  <td> DoD ID#: </td>
                  <td> 1111111111 </td>
                </tr>
                <tr>
                  <td> Current Duty Station: </td>
                  <td> Joint Base San Antonio (JBSA) </td>
                </tr>
              </table>

              <table>
                <tr>
                  <th>Orders Edit</th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>
                    Michael Owen Jones{' '}
                    {currentServiceMember && currentServiceMember.first_name}
                  </td>
                </tr>
                <tr>
                  <td> Branch: </td>
                  <td> Air Force </td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td> E-5 </td>
                </tr>
                <tr>
                  <td> DoD ID#: </td>
                  <td> 1111111111 </td>
                </tr>
                <tr>
                  <td> Current Duty Station: </td>
                  <td> Joint Base San Antonio (JBSA) </td>
                </tr>
              </table>
            </div>

            <div className="usa-width-one-half review-section">
              <table>
                <tr>
                  <th>Contact Info Edit</th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>
                    Michael Owen Jones{' '}
                    {currentServiceMember && currentServiceMember.first_name}
                  </td>
                </tr>
                <tr>
                  <td> Branch: </td>
                  <td> Air Force </td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td> E-5 </td>
                </tr>
                <tr>
                  <td> DoD ID#: </td>
                  <td> 1111111111 </td>
                </tr>
                <tr>
                  <td> Current Duty Station: </td>
                  <td> Joint Base San Antonio (JBSA) </td>
                </tr>
              </table>

              <table>
                <tr>
                  <th>
                    Backup Contact Info{' '}
                    <span className="align-right">
                      <a href="about:blank">Edit</a>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>
                    {currentBackupContacts[0] && currentBackupContacts[0].email}
                  </td>
                </tr>
                <tr>
                  <td> Branch: </td>
                  <td> Air Force </td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td> E-5 </td>
                </tr>
                <tr>
                  <td> DoD ID#: </td>
                  <td> 1111111111 </td>
                </tr>
                <tr>
                  <td> Current Duty Station: </td>
                  <td> Joint Base San Antonio (JBSA) </td>
                </tr>
              </table>
            </div>
          </div>
          <p>{currentPpm && currentPpm.estimated_incentive}</p>
        </Fragment>
      </WizardPage>
    );
  }
}

Review.propTypes = {
  currentServiceMember: PropTypes.object,
  currentPpm: PropTypes.object,
  currentBackupContacts: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPpm }, dispatch);
}

function mapStateToProps(state) {
  const props = {
    ...state.serviceMember,
    ...state.ppm,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Review);

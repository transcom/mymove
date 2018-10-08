import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PropTypes from 'prop-types';
import { isOfficeSite } from 'shared/constants.js';

import PreApprovalTable from 'shared/PreApprovalRequest/PreApprovalTable.jsx';
import Creator from 'shared/PreApprovalRequest/Creator';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

class PreApprovalPanel extends Component {
  onSubmit = values => {
    return new Promise(function(resolve, reject) {
      // do a thing, possibly async, then…
      setTimeout(function() {
        console.log('onSubmit async', values);
        resolve('success');
      }, 50);
    });
  };
  onEdit = () => {
    console.log('onEdit hit');
  };
  onDelete = () => {
    console.log('onDelete hit');
  };
  onApproval = () => {
    console.log('onApproval hit');
  };

  render() {
    return (
      <div>
        <BasicPanel title={'Pre-Approval Requests'}>
          <PreApprovalTable
            shipment_accessorials={this.props.shipment_accessorials}
            isActionable={true}
            onEdit={this.onEdit}
            onDelete={this.onDelete}
            onApproval={isOfficeSite ? this.onApproval : null}
          />
          <Creator
            accessorials={this.props.accessorials}
            savePreApprovalRequest={this.onSubmit}
          />
        </BasicPanel>
      </div>
    );
  }
}

PreApprovalPanel.propTypes = {
  shipment_accessorials: PropTypes.array,
  accessorials: PropTypes.array,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);

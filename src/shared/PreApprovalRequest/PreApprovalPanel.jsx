import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PropTypes from 'prop-types';
import PreApprovalRequestForm, {
  formName as PreApprovalRequestFormName,
} from 'shared/PreApprovalRequestForm';
import { submit, isValid, isSubmitting } from 'redux-form';
import TableList from 'shared/PreApprovalRequest/TableList.jsx';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

class PreApprovalPanel extends Component {
  // TODO - should onSubmit be defined in the Creator and Editor components accordingly?
  onSubmit = values => {
    console.log('onSubmit', values);
  };
  onEdit = () => {
    console.log('onEdit');
  };
  onDelete = () => {
    console.log('onDelete');
  };
  onApproval = () => {
    console.log('onApproval hit');
  };

  render() {
    return (
      <div>
        <BasicPanel title={'PreApproval Requests'}>
          <TableList
            shipment_accessorials={this.props.shipment_accessorials}
            // TODO: Check if office or tsp app - set true for office
            isActionable={true}
            onEdit={this.onEdit}
            onDelete={this.onDelete}
            onApproval={this.onApproval}
          />
          <PreApprovalRequestForm
            accessorials={this.props.accessorials}
            ref={form => (this.formReference = form)}
            onSubmit={this.onSubmit}
          />
          <button
            disabled={!this.props.formEnabled}
            onClick={this.props.submitForm}
          >
            Submit
          </button>
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
  return {
    formEnabled:
      isValid(PreApprovalRequestFormName)(state) &&
      !isSubmitting(PreApprovalRequestFormName)(state),
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(PreApprovalRequestFormName),
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);

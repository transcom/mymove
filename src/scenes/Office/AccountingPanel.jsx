import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import { updateAccounting, getAccounting } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import { EditablePanel, EditableTextField } from 'shared/EditablePanel';

const formName = 'office_move_info_accounting';
class AccountingPanel extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isEditable: false,
    };
    this.toggleEditable = this.toggleEditable.bind(this);
  }

  componentDidMount() {
    this.props.getAccounting(this.props.moveID);
  }

  toggleEditable() {
    this.setState({ isEditable: !this.state.isEditable });
  }

  render() {
    return (
      <EditablePanel
        title="Accounting"
        editableComponent={AccountingPanelEditable}
        displayComponent={AccountingPanelDisplay}
        isEditable={this.state.isEditable}
        toggleEditable={this.toggleEditable}
      />
    );
  }
}

// TODO add proptypes

const AccountingPanelDisplay = props => {
  return (
    <div>
      <p>TAC Value</p>
    </div>
  );
};

const AccountingPanelEditable = props => {
  const schema = props.schema;

  return (
    <div>
      <SwaggerField fieldName="dept_indicator" swagger={schema} required />
      <SwaggerField fieldName="tac" swagger={schema} required />
    </div>
  );
};

AccountingPanel = reduxForm({ form: formName })(AccountingPanel);

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateAccounting, getAccounting }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: {},
    formData: state.form[formName],
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.PatchAccounting;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(AccountingPanel);

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import { updateAccounting, getAccounting } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import { EditablePanel } from 'shared/EditablePanel';

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
    const displayContent = () => {
      return (
        <div>
          <p>TAC Value</p>
        </div>
      );
    };

    const editableContent = () => {
      const { schema } = this.props;
      return (
        <div>
          <SwaggerField fieldName="dept_indicator" swagger={schema} required />
          <SwaggerField fieldName="tac" swagger={schema} required />
        </div>
      );
    };

    return (
      <EditablePanel
        title="Accounting"
        editableContent={editableContent}
        displayContent={displayContent}
        isEditable={this.state.isEditable}
        toggleEditable={this.toggleEditable}
      />
    );
  }
}

AccountingPanel.propTypes = {
  schema: PropTypes.object.isRequired,
};

const formName = 'office_move_info_accounting';
AccountingPanel = reduxForm({ form: formName })(AccountingPanel);

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      updateAccounting,
      getAccounting,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    schema: {},
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.PatchAccounting;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(AccountingPanel);

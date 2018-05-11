import { get } from 'lodash';
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import { updateAccounting, loadAccounting } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Alert from 'shared/Alert';
import { EditablePanel } from 'shared/EditablePanel';

const PanelField = props => {
  return (
    <div className="panel-field">
      <span className="field-title">{props.title}</span>
      <span className="field-value">{props.value}</span>
    </div>
  );
};

class AccountingPanel extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isEditable: false,
    };
  }

  componentDidMount() {
    this.props.loadAccounting(this.props.moveId);
  }

  save = () => {
    this.props.updateAccounting(this.props.moveId, this.props.formData.values);
    this.toggleEditable();
  };

  toggleEditable = () => {
    this.setState({
      isEditable: !this.state.isEditable,
    });
  };

  render() {
    const displayContent = () => {
      const values = this.props.accounting || {};
      return (
        <React.Fragment>
          <div className="editable-panel-column">
            <PanelField title="Dept. indicator" value={values.dept_indicator} />
          </div>
          <div className="editable-panel-column">
            <PanelField title="TAC" value={values.tac} />
          </div>
        </React.Fragment>
      );
    };

    const editableContent = () => {
      const { schema } = this.props;
      return (
        <React.Fragment>
          <div className="editable-panel-column">
            <SwaggerField
              fieldName="dept_indicator"
              swagger={schema}
              required
            />
          </div>
          <div className="editable-panel-column">
            <SwaggerField fieldName="tac" swagger={schema} required />
          </div>
        </React.Fragment>
      );
    };

    return (
      <React.Fragment>
        {this.props.hasError && (
          <Alert type="error" heading="An error occurred">
            There was an error: <em>{this.props.errorMessage}</em>.
          </Alert>
        )}
        <EditablePanel
          title="Accounting"
          editableContent={editableContent}
          displayContent={displayContent}
          onSave={this.save}
          toggleEditable={this.toggleEditable}
          isEditable={this.state.isEditable || this.props.isUpdating}
        />
      </React.Fragment>
    );
  }
}

AccountingPanel.propTypes = {
  schema: PropTypes.object.isRequired,
  moveId: PropTypes.string.isRequired,
};

const formName = 'office_move_info_accounting';
AccountingPanel = reduxForm({ form: formName })(AccountingPanel);

function mapStateToProps(state) {
  return {
    schema: get(state, 'swagger.spec.definitions.PatchAccounting', {}),
    hasError:
      state.officeAccounting.hasLoadError ||
      state.officeAccounting.hasUpdateError,
    errorMessage: state.officeAccounting.error,
    formData: state.form[formName],
    initialValues: state.officeAccounting.accounting,
    ...state.officeAccounting,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      updateAccounting,
      loadAccounting,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(AccountingPanel);

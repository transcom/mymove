import { get } from 'lodash';
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import { updateAccounting, loadAccounting } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Alert from 'shared/Alert';
import { EditablePanel, PanelField } from 'shared/EditablePanel';

const AccountingDisplay = props => {
  const { values } = props;
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

const AccountingEdit = props => {
  const { schema } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="dept_indicator" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="tac" swagger={schema} required />
      </div>
    </React.Fragment>
  );
};

function editablePanel(DisplayComponent, EditComponent, reducerName) {
  let Wrapper = class extends Component {
    constructor(props) {
      super(props);
      this.state = {
        isEditable: false,
      };
    }

    componentDidMount() {
      this.props.load(this.props.moveId);
    }

    save = () => {
      this.props.update(this.props.moveId, this.props.formData.values);
      this.toggleEditable();
    };

    toggleEditable = () => {
      this.setState({
        isEditable: !this.state.isEditable,
      });
    };

    render() {
      const Content = this.state.isEditable ? EditComponent : DisplayComponent;

      return (
        <React.Fragment>
          {this.props.hasError && (
            <Alert type="error" heading="An error occurred">
              There was an error: <em>{this.props.errorMessage}</em>.
            </Alert>
          )}
          <EditablePanel
            title="Accounting"
            onSave={this.save}
            toggleEditable={this.toggleEditable}
            isEditable={this.state.isEditable || this.props.isUpdating}
          >
            <Content
              values={this.props.displayValues}
              schema={this.props.schema}
            />
          </EditablePanel>
        </React.Fragment>
      );
    }
  };

  Wrapper.propTypes = {
    schema: PropTypes.object.isRequired,
    displayValues: PropTypes.object.isRequired,
    load: PropTypes.func.isRequired,
    update: PropTypes.func.isRequired,
    moveId: PropTypes.string.isRequired,
  };

  const formName = `office_move_info_${reducerName}`;
  Wrapper = reduxForm({ form: formName })(Wrapper);

  function mapStateToProps(globalState) {
    const state = get(globalState, reducerName, {});
    return {
      // Wrapper
      schema: get(globalState, 'swagger.spec.definitions.PatchAccounting', {}),
      hasError: state.hasLoadError || state.hasUpdateError,
      errorMessage: state.error,
      displayValues: state.accounting || {},
      isUpdating: state.isUpdating,

      // reduxForm
      formData: globalState.form[formName],
      initialValues: state.accounting,
    };
  }

  function mapDispatchToProps(dispatch) {
    return bindActionCreators(
      {
        update: updateAccounting,
        load: loadAccounting,
      },
      dispatch,
    );
  }

  return connect(mapStateToProps, mapDispatchToProps)(Wrapper);
}

export default editablePanel(
  AccountingDisplay,
  AccountingEdit,
  'officeAccounting',
);

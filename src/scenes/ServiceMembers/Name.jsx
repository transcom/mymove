import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import { updateServiceMember } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import SectionWrapper from 'components/Customer/SectionWrapper';

const subsetOfFields = ['first_name', 'middle_name', 'last_name', 'suffix'];
const formName = 'service_member_name';
const NameWizardForm = reduxifyWizardForm(formName);

export class Name extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      return this.props.updateServiceMember(patch);
    }
  };

  render() {
    const { pages, pageKey, error, currentServiceMember } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <NameWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <h1>Name</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <SwaggerField fieldName="first_name" swagger={this.props.schema} required />
          </div>
          <SwaggerField fieldName="middle_name" swagger={this.props.schema} />
          <SwaggerField fieldName="last_name" swagger={this.props.schema} required />
          <SwaggerField fieldName="suffix" swagger={this.props.schema} />
        </SectionWrapper>
      </NameWizardForm>
    );
  }
}
Name.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  return {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(Name);

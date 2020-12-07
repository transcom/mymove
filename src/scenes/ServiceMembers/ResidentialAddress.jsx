import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { ValidateZipRateData } from 'shared/api';
import AddressForm from 'shared/AddressForm';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

import SectionWrapper from 'components/Customer/SectionWrapper';

const UnsupportedZipCodeErrorMsg =
  'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.';

async function asyncValidate(values) {
  const { postal_code } = values;
  const responseBody = await ValidateZipRateData(postal_code, 'origin');
  if (!responseBody.valid) {
    // eslint-disable-next-line no-throw-literal
    throw { postal_code: UnsupportedZipCodeErrorMsg };
  }
}

const formName = 'service_member_residential_address';
const ResidentalWizardForm = reduxifyWizardForm(formName, null, asyncValidate, ['postal_code']);

export class ResidentialAddress extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  handleSubmit = () => {
    const { values, currentServiceMember, updateServiceMember } = this.props;

    const payload = {
      id: currentServiceMember.id,
      residential_address: values,
    };

    return patchServiceMember(payload)
      .then((response) => {
        updateServiceMember(response);
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        this.setState({
          errorMessage,
        });
      });
  };

  render() {
    const { pages, pageKey, currentServiceMember } = this.props;
    const { errorMessage } = this.state;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = get(currentServiceMember, 'residential_address');
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <ResidentalWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={errorMessage}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <h1>Current residence</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <AddressForm schema={this.props.schema} />
          </div>
        </SectionWrapper>
      </ResidentalWizardForm>
    );
  }
}
ResidentialAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    schema: get(state, 'swaggerInternal.spec.definitions.Address', {}),
    values: getFormValues(formName)(state),
    currentServiceMember: serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ResidentialAddress);

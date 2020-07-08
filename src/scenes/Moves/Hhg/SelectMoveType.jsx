import React, { Component } from 'react';
import RadioButton from 'shared/RadioButton';
import { SHIPMENT_TYPE } from 'shared/constants';
import WizardPage from 'shared/WizardPage';
import { no_op } from 'shared/utils';

export class SelectMoveType extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      moveType: SHIPMENT_TYPE.PPM,
    };
  }

  handleRadioChange = (event) => {
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  render() {
    const { pages, pageKey, error } = this.props;

    return (
      <WizardPage handleSubmit={no_op} pageList={pages} pageKey={pageKey} error={error}>
        <div className="grid-container usa-prose">
          <div className="usa-grid">
            <div className="grid-row grid-gap">
              <h1 className="sm-heading">How do you want to move your belongings?</h1>
              <div className="grid-col-9 desktop:grid-col-12">
                <RadioButton
                  inputClassName="inline_radio"
                  labelClassName="inline_radio"
                  label="Arrange it all yourself"
                  value={SHIPMENT_TYPE.PPM}
                  name="moveType"
                  checked={this.state.moveType === SHIPMENT_TYPE.PPM}
                  onChange={this.handleRadioChange}
                />
              </div>
            </div>
            <div className="grid-row grid-gap">
              <div className="grid-col-9 desktop:grid-col-12">
                <RadioButton
                  inputClassName="inline_radio"
                  labelClassName="inline_radio"
                  label="Have professionals pack and move it all"
                  value={SHIPMENT_TYPE.HHG}
                  name="moveType"
                  // TODO: uncomment when we have more HHG pages
                  // checked={this.state.moveType === SHIPMENT_TYPE.HHG}
                  disabled={true}
                  onChange={this.handleRadioChange}
                />
              </div>
            </div>
          </div>
        </div>
      </WizardPage>
    );
  }
}

export default SelectMoveType;

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import RadioButton from 'shared/RadioButton';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { MOVE_TYPES } from 'shared/constants';

export class SelectMoveType extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      moveType: MOVE_TYPES.PPM,
    };
  }

  handleRadioChange = (event) => {
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  render() {
    return (
      <div className="grid-container usa-prose">
        <div className="usa-grid">
          <div className="grid-row grid-gap">
            <h1 className="sm-heading">How do you want to move your belongings?</h1>
            <div className="grid-col-9 desktop:grid-col-12">
              <RadioButton
                inputClassName="inline_radio"
                labelClassName="inline_radio"
                label="Arrange it all yourself"
                value={MOVE_TYPES.PPM}
                name="moveType"
                checked={this.state.moveType === MOVE_TYPES.PPM}
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
                value={MOVE_TYPES.HHG}
                name="moveType"
                checked={this.state.moveType === MOVE_TYPES.HHG}
                onChange={this.handleRadioChange}
              />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const move = selectActiveOrLatestMove(state);
  const props = {
    moveId: get(move, 'id'),
    selectedMoveType: get(move, 'selected_move_type'),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);

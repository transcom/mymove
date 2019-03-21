import React from 'react';
import ComboButton from 'shared/ComboButton/index.jsx';
import { DropDown, DropDownItem } from 'shared/ComboButton/dropdown';
import ToolTip from 'shared/ToolTip';

const StyleGuide = () => (
  <div style={{ 'margin-left': '20px' }}>
    <h1>Wizard Styles</h1>
    <hr />
    <h2>This is a H2</h2>
    <h3>This is a H3</h3>
    <ComboButton buttonText="Approve" disabled={false}>
      <DropDown>
        <DropDownItem disabled={false} value="Enabled Menu Item" />
        <DropDownItem disabled={true} value="Disabled Menu Item" />
      </DropDown>
    </ComboButton>
    <span style={{ 'margin-left': '30px' }}>
      <ToolTip textStyle="tooltiptext-medium" disabled={false} toolTipText="Tooltip text">
        <ComboButton disabled={true} buttonText="Approve" />
      </ToolTip>
    </span>
  </div>
);

export default StyleGuide;

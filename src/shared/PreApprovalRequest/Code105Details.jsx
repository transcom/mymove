import React from 'react';
import { DimensionsField } from '../JsonSchemaForm/DimensionsField';

export const Code105Details = props => {
  return (
    <div>
      <DimensionsField fieldName="item_dimensions" swagger={props.swagger} labelText="Item Dimensions (inches)" />
      <DimensionsField fieldName="crate_dimensions" swagger={props.swagger} labelText="Crate Dimensions (inches)" />
      <div className="bq-explanation">
        <p>Crate can only exceed item size by:</p>
        <ul>
          <li>
            <em>Internal crate</em>: Up to 3" larger
          </li>
          <li>
            <em>External crate</em>: Up to 5" larger
          </li>
        </ul>
      </div>
    </div>
  );
};

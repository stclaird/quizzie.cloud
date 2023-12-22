import React from 'react';
export default function SubmitButton(props) {
    return (
    <button
        className="form-submit-button"
        onClick={
            () => props.submitAnswer(props.id)
        }
        >Submit
    </button>
    );
  }
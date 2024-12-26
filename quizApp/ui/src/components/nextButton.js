import React from 'react';
export default function NextButton(props) {
    function isLast() {
        let lastQuestion = props.NumQuestions - 1
        if (lastQuestion === props.QuestionIdx) {
            return true
        }
        return false
    }

    return (
   <button
        className="form-button-next"
        disabled={isLast()}
        onClick={() => props.NextQuestion()
       }
        >Next Question &gt;
    </button>
    );
  }
import React, { useState , useEffect} from 'react';
export default function Category(props) {
    const [categories, setCategories] = useState([]);
    const fetchCategories = async () => {
        try {
          const response = await fetch(props.categoriesUrl, {
            method: 'GET'
          });
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }

          const data = await response.json();
          setCategories(data);
        } catch (error) {
          console.error('Error:', error);
        }
      }

      useEffect(() => {
        fetchCategories();
        }, []);

    return (
            <div >
            {categories.map((category, index) => (
              <div className="categories" key={index}>
                <h2 id="list-heading">{category.category}</h2>

                {category.subcategories.map((subCat, index) => (
                  <div key={index}>
                    <div key={subCat.subcategoryname} className="btn-group">
                        <button
                            type="button"
                            key={subCat.subcategoryname}
                            className="btn btn__primary"
                            onClick={() => props.startQuestions(category.category, subCat.subcategoryname, subCat.urlprefix)}
                            >
                        {subCat.subcategoryname}
                        </button>
                    </div>
                  </div>
                ))}
              </div>
            ))}
          </div>
        );
  }